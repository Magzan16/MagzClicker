#!/usr/bin/env python3
"""Embed an ICO file as the application icon of a PE/PE32+ executable.

Pure Python (standard library only). It appends a .rsrc section containing
RT_ICON and RT_GROUP_ICON resources. Designed for MagzClicker's reproducible
Windows build; it does not alter program code.
"""
from __future__ import annotations

import argparse
import struct
from pathlib import Path

RT_ICON = 3
RT_GROUP_ICON = 14
LANG_EN_US = 1033
SECTION_READ_INITIALIZED_DATA = 0x40000040


def align(value: int, alignment: int) -> int:
    return (value + alignment - 1) // alignment * alignment


def u16(buf: bytes | bytearray, off: int) -> int:
    return struct.unpack_from('<H', buf, off)[0]


def u32(buf: bytes | bytearray, off: int) -> int:
    return struct.unpack_from('<I', buf, off)[0]


def parse_ico(path: Path):
    data = path.read_bytes()
    if len(data) < 6:
        raise ValueError('ICO file is too small')
    reserved, typ, count = struct.unpack_from('<HHH', data, 0)
    if reserved != 0 or typ != 1 or count < 1:
        raise ValueError('Not a valid icon (.ico) file')
    if len(data) < 6 + count * 16:
        raise ValueError('Truncated ICO directory')

    frames = []
    for i in range(count):
        off = 6 + i * 16
        width, height, colors, reserved_b, planes, bitcount, size, image_off = struct.unpack_from('<BBBBHHII', data, off)
        if image_off + size > len(data):
            raise ValueError('Truncated ICO image data')
        frames.append({
            'width': width,
            'height': height,
            'colors': colors,
            'reserved': reserved_b,
            'planes': planes,
            'bitcount': bitcount,
            'data': data[image_off:image_off + size],
        })
    return frames


def dir_header(id_count: int) -> bytes:
    return struct.pack('<IIHHHH', 0, 0, 0, 0, 0, id_count)


def dir_entry(resource_id: int, target_off: int, is_directory: bool) -> bytes:
    return struct.pack('<II', resource_id, target_off | (0x80000000 if is_directory else 0))


def build_rsrc(frames, section_rva: int) -> bytes:
    n = len(frames)

    # Directory layout: root -> type -> resource ID -> language -> data entry.
    root_off = 0
    root_size = 16 + 2 * 8
    icon_type_off = root_off + root_size
    icon_type_size = 16 + n * 8
    group_type_off = icon_type_off + icon_type_size
    group_type_size = 16 + 1 * 8

    icon_lang_offs = []
    cursor = group_type_off + group_type_size
    for _ in range(n):
        icon_lang_offs.append(cursor)
        cursor += 16 + 1 * 8
    group_lang_off = cursor
    cursor += 16 + 1 * 8

    icon_data_entry_offs = []
    for _ in range(n):
        icon_data_entry_offs.append(cursor)
        cursor += 16
    group_data_entry_off = cursor
    cursor += 16

    payload_cursor = align(cursor, 4)
    icon_payload_offs = []
    for frame in frames:
        icon_payload_offs.append(payload_cursor)
        payload_cursor = align(payload_cursor + len(frame['data']), 4)

    group_payload_off = payload_cursor
    group_data = bytearray(struct.pack('<HHH', 0, 1, n))
    for idx, frame in enumerate(frames, start=1):
        group_data += struct.pack(
            '<BBBBHHIH',
            frame['width'], frame['height'], frame['colors'], frame['reserved'],
            frame['planes'], frame['bitcount'], len(frame['data']), idx,
        )
    total_size = align(group_payload_off + len(group_data), 4)
    out = bytearray(total_size)

    def put(off: int, blob: bytes | bytearray):
        out[off:off + len(blob)] = blob

    # Root directory: numeric type IDs must be sorted.
    put(root_off, dir_header(2))
    put(root_off + 16, dir_entry(RT_ICON, icon_type_off, True))
    put(root_off + 24, dir_entry(RT_GROUP_ICON, group_type_off, True))

    put(icon_type_off, dir_header(n))
    for i in range(n):
        put(icon_type_off + 16 + i * 8, dir_entry(i + 1, icon_lang_offs[i], True))

    put(group_type_off, dir_header(1))
    put(group_type_off + 16, dir_entry(1, group_lang_off, True))

    for i in range(n):
        put(icon_lang_offs[i], dir_header(1))
        put(icon_lang_offs[i] + 16, dir_entry(LANG_EN_US, icon_data_entry_offs[i], False))

    put(group_lang_off, dir_header(1))
    put(group_lang_off + 16, dir_entry(LANG_EN_US, group_data_entry_off, False))

    for i, frame in enumerate(frames):
        put(icon_data_entry_offs[i], struct.pack(
            '<IIII', section_rva + icon_payload_offs[i], len(frame['data']), 0, 0
        ))
        put(icon_payload_offs[i], frame['data'])

    put(group_data_entry_off, struct.pack(
        '<IIII', section_rva + group_payload_off, len(group_data), 0, 0
    ))
    put(group_payload_off, group_data)
    return bytes(out)


def patch_pe(exe_path: Path, ico_path: Path, output_path: Path | None):
    data = bytearray(exe_path.read_bytes())
    if data[:2] != b'MZ':
        raise ValueError('Input is not a PE executable (missing MZ header)')

    pe_off = u32(data, 0x3C)
    if data[pe_off:pe_off + 4] != b'PE\0\0':
        raise ValueError('Input is not a valid PE executable')

    coff = pe_off + 4
    num_sections = u16(data, coff + 2)
    size_opt = u16(data, coff + 16)
    opt = coff + 20
    magic = u16(data, opt)
    if magic not in (0x10B, 0x20B):
        raise ValueError(f'Unsupported PE optional-header magic: 0x{magic:04x}')

    section_alignment = u32(data, opt + 32)
    file_alignment = u32(data, opt + 36)
    size_headers = u32(data, opt + 60)
    size_image_off = opt + 56
    size_init_data_off = opt + 8

    data_dir_off = opt + (96 if magic == 0x10B else 112)
    resource_dir_off = data_dir_off + 2 * 8
    old_rsrc_rva = u32(data, resource_dir_off)
    old_rsrc_size = u32(data, resource_dir_off + 4)
    if old_rsrc_rva or old_rsrc_size:
        raise ValueError('Executable already contains a resource directory; refusing to overwrite it')

    sec_table = opt + size_opt
    new_sec_hdr_off = sec_table + num_sections * 40
    if new_sec_hdr_off + 40 > size_headers:
        raise ValueError('Not enough PE header space to add a .rsrc section')

    max_end_rva = 0
    max_end_raw = 0
    for i in range(num_sections):
        sh = sec_table + i * 40
        vsize = u32(data, sh + 8)
        vaddr = u32(data, sh + 12)
        raw_size = u32(data, sh + 16)
        raw_ptr = u32(data, sh + 20)
        max_end_rva = max(max_end_rva, vaddr + max(vsize, raw_size))
        max_end_raw = max(max_end_raw, raw_ptr + raw_size)

    new_rva = align(max_end_rva, section_alignment)
    raw_ptr = align(max(len(data), max_end_raw), file_alignment)
    frames = parse_ico(ico_path)
    rsrc = build_rsrc(frames, new_rva)
    raw_size = align(len(rsrc), file_alignment)

    # New section header.
    section_header = struct.pack(
        '<8sIIIIIIHHI',
        b'.rsrc\0\0\0', len(rsrc), new_rva, raw_size, raw_ptr,
        0, 0, 0, 0, SECTION_READ_INITIALIZED_DATA,
    )
    data[new_sec_hdr_off:new_sec_hdr_off + 40] = section_header

    # COFF/optional header updates.
    struct.pack_into('<H', data, coff + 2, num_sections + 1)
    old_init = u32(data, size_init_data_off)
    struct.pack_into('<I', data, size_init_data_off, old_init + raw_size)
    struct.pack_into('<I', data, size_image_off, align(new_rva + len(rsrc), section_alignment))
    struct.pack_into('<II', data, resource_dir_off, new_rva, len(rsrc))

    # Append the new section without shifting any existing PE data.
    if len(data) < raw_ptr:
        data.extend(b'\0' * (raw_ptr - len(data)))
    data.extend(rsrc)
    if len(rsrc) < raw_size:
        data.extend(b'\0' * (raw_size - len(rsrc)))

    out = output_path or exe_path
    out.write_bytes(data)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument('exe', type=Path)
    ap.add_argument('ico', type=Path)
    ap.add_argument('-o', '--output', type=Path)
    args = ap.parse_args()
    patch_pe(args.exe, args.ico, args.output)


if __name__ == '__main__':
    main()
