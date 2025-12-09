# openlist-bed

An image bed service based on [openlist](https://github.com/OpenListTeam/OpenList) 
and [libvips](https://github.com/libvips/libvips).


## Features
- Supports multiple image formats including JPEG, PNG, WebP, 
Heif (only read, will transfer to JPEG or Webp), GIF.
- Automatic compression and transfer to WebP format for smaller file sizes.
- Auto naming of images using UNIX and image width and height.

## Installation

check [docker-compost.yml](./docker-compose.yml).

## Configuration

| env key       | example value           | description                               |
|---------------|-------------------------|-------------------------------------------|
| API_KEY       | RANDOM_KEY              | the api-key to access api endpoint        |
| HOST          | https://as.example.com  | to generate output image url              |
| OPENLIST_HOST | https://pan.example.com | openlist service address                  |
| OPENLIST_PATH | /path/to                | the path to store images at openlist      |
| AUTO_WEBP     | false                   | to enable auto check browser webp support |
