# AssetGen

AssetGen converts PNGs to asset catalogs. This program as a very narrow focus, see assumptions below.

## Assumptions

- You will run this on a Mac.
- If an images path contains `.imageset/` or `.xcassets/`, it will be ignored.
- Only files ending in `.png` are processed. This is case sensitive!
- The program assumes your files have the follow suffixes `.png`, `@2x.png`, and `@3x.png`. These are case sensitive!
- If the `.imageset` folder already exists, the file set will be skipped.

## Installation

### Download binary

You can download the binary directly from [here](https://github.com/unrolled/assetgen/releases/download/v1.0.0/assetgen). Once the binary has been downloaded, you can move it to `/usr/local/bin` and give it execution privileges:

    mv assetgen /usr/local/bin/assetgen
    chmod +x /usr/local/bin/assetgen

### Installing with `go get`

Ensure [Go](https://golang.org/dl/) is installed, and run the following:

    go get github.com/unrolled/assetgen

## Usage

If you are in your images folder, you can issue the command like this:

    assetgen .

Otherwise, you can specify any valid directory path:

    assetgen ~/code/myapp/resources/
