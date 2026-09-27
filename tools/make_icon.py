#!/usr/bin/env python3
"""Generate the MagzClicker ICO file using only the Python standard library."""

from __future__ import annotations

import argparse
import struct
from pathlib import Path


def canvas(size: int):
    return [[[0, 0, 0, 0] for _ in range(size)] for _ in range(size)]


def fill(px, x0, y0, x1, y1, color):
    h = len(px)
    w = len(px[0])
    x0, y0 = max(0, x0), max(0, y0)
    x1, y1 = min(w, x1), min(h, y1)
    for y in range(y0, y1):
        row = px[y]
        for x in range(x0, x1):
            row[x] = list(color)


def draw_icon(size: int):
    px = canvas(size)
    u = size / 64.0

    def r(x0, y0, x1, y1, color):
        fill(px, round(x0 * u), round(y0 * u), round(x1 * u), round(y1 * u), color)

    # Transparent margin and dark outline.
    r(5, 6, 59, 58, (65, 20, 10, 255))
    r(7, 8, 57, 56, (185, 35, 20, 255))

    # Red/orange striped shell.
    stripes = [
        (7, 15, (225, 48, 24, 255)),
        (15, 22, (133, 31, 23, 255)),
        (22, 31, (242, 67, 25, 255)),
        (31, 39, (150, 32, 20, 255)),
        (39, 48, (235, 55, 22, 255)),
        (48, 57, (132, 29, 18, 255)),
    ]
    for x0, x1, c in stripes:
        r(x0, 8, x1, 56, c)

    # Top and bottom highlights.
    r(8, 9, 56, 12, (255, 108, 31, 255))
    r(8, 52, 56, 55, (110, 23, 18, 255))

    # TNT label band.
    r(7, 24, 57, 42, (230, 231, 228, 255))
    r(7, 24, 57, 27, (250, 250, 248, 255))
    r(7, 39, 57, 42, (184, 187, 184, 255))

    # Pixel letters on a 3 x 5 grid.
    patterns = {
        "T": ["111", "010", "010", "010", "010"],
        "N": ["101", "111", "111", "111", "101"],
    }
    letter_w = 3
    letter_h = 5
    cell = 3.0
    gap = 2.0
    total = (letter_w * cell) * 3 + gap * 2
    start_x = 32 - total / 2
    start_y = 26.5
    for idx, ch in enumerate("TNT"):
        pattern = patterns[ch]
        ox = start_x + idx * (letter_w * cell + gap)
        for yy in range(letter_h):
            for xx in range(letter_w):
                if pattern[yy][xx] == "1":
                    r(ox + xx * cell, start_y + yy * cell,
                      ox + (xx + 1) * cell, start_y + (yy + 1) * cell,
                      (28, 28, 28, 255))

    # Small corner shine.
    r(9, 10, 12, 13, (255, 155, 45, 255))
    return px


def dib_frame(px) -> bytes:
    size = len(px)
    # ICO DIB height includes XOR + AND mask heights.
    header = struct.pack(
        "<IIIHHIIIIII",
        40, size, size * 2, 1, 32, 0,
        size * size * 4, 0, 0, 0, 0,
    )
    xor = bytearray()
    for y in range(size - 1, -1, -1):
        for r, g, b, a in px[y]:
            xor.extend((b, g, r, a))
    row_bytes = ((size + 31) // 32) * 4
    and_mask = b"\x00" * (row_bytes * size)
    return header + bytes(xor) + and_mask


def write_ico(path: Path):
    sizes = (16, 32, 48, 64)
    frames = [(s, dib_frame(draw_icon(s))) for s in sizes]
    header = struct.pack("<HHH", 0, 1, len(frames))
    offset = 6 + 16 * len(frames)
    directory = bytearray()
    payload = bytearray()
    for size, frame in frames:
        wh = 0 if size == 256 else size
        directory.extend(struct.pack("<BBBBHHII", wh, wh, 0, 0, 1, 32, len(frame), offset))
        payload.extend(frame)
        offset += len(frame)
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_bytes(header + directory + payload)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("output", type=Path)
    args = parser.parse_args()
    write_ico(args.output)


if __name__ == "__main__":
    main()
