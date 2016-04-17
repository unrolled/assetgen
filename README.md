# AssetGen

AssetGen converts PNGs to asset catalogs. This program as a very narrow focus, see assumptions below.

## Assumptions

- You will run this on a Mac.
- If an images path contains `.imageset/` or `.xcassets/`, it will be ignored.
- Only files ending in `.png` are processed. This is case sensitive!
- The program assumes your files have the follow suffixes `.png`, `@2x.png`, and `@3x.png`. These are case sensitive!
- If the `.imageset` folder already exists, the file set will be skipped.

## Installation

Ensure [Go](https://golang.org/dl/) is installed, and run the following:

    go get github.com/unrolled/assetgen

## Usage

If you are in your images folder, you can issue the command like this:

    assetgen .

Otherwise, you can specify any valid directory path:

    assetgen ~/code/myapp/resources/
