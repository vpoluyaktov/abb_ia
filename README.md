# Audiobook Builder (Internet Archive version)

## Description

The Internet Archive site offers a vast collection of free "old-time radio" shows, audiobooks, and lectures that can be downloaded. While you can listen to them on your web browser, it can be inconvenient especially if you want to listen to them on your mobile device. Typically, these shows are provided as individual .mp3 files, requiring you to download them all, create a playlist, and remember your last listened file and position.

To make this process easier, I developed Audiobook Builder. With this app, all you need is the name of a show or book, or a direct link on archive.org. It will download the .mp3 files for the book, re-encode them with the same bit rate, generate a list of chapters (which can be edited during the process), and ultimately create an audiobook in .m4b format.


![Audiobook Builder in action](https://github.com/vpoluyaktov/abb_ia/blob/master/assets/abb_ia.gif)

Here is what the newly created book looks like in the **Audiobookshelf Web** browser and the **Audiobookshelf iOS app**:
![Created book in Audiobookshelf browser](https://github.com/vpoluyaktov/abb_ia/blob/master/assets/audiobookshelf_browser.png)
![Created book in Audiobookshelf IOS app](https://github.com/vpoluyaktov/abb_ia/blob/master/assets/audiobookshelf_ios.png)

## Features

- Download a set of single .mp3 files from [archive.org](https://archive.org)
- Create an audiobook in .m4b format
- Re-encode mp3 files to the same bit rate, if necessary.
- Modify audiobook metadata obtained from [archive.org](https://archive.org), including book title, author, series, genre, and art cover
- Upload created audiobook to specified folder using [audiobookshelf compatible directory structure](https://www.audiobookshelf.org/docs/#book-directory-structure)
- Upload created audiobook to remote Audiobookshelf server using [Audiobookshelf API](https://api.audiobookshelf.org/)

## Integrations

Audiobook Builder seamlessly integrates with **Audiobookshelf server** (https://www.audiobookshelf.org). This integration allows you to upload the created audiobooks directly to the Audiobook Shelf server for convenient listening.

## How to Start

1. Make sure you have ffmpeg and ffprobe command line utilities installed. If not, install them first.
2. Download the ready-to-run binary files from the [Github Releases page](https://github.com/vpoluyaktov/abb_ia/releases) for your target platform.
3. Open a terminal and navigate to the directory where the binary file is located.
4. Run the binary file by executing the command `./abb_ia`. The TUI interface will appear.

## Installation Instructions

To use Audiobook Builder, you need to have the following command line utilities installed:

- **ffmpeg** (used for audio manipulation)
- **ffprobe** (used for retrieving audio metadata)

Make sure these utilities are properly installed and available in your system's PATH before running `abb_ia`.

To install Audiobook Builder (abb_ia) on your system, follow these steps:

1. Download the ready-to-run binary file for your target platform from the [Github Releases page](https://github.com/vpoluyaktov/abb_ia/releases).
2. Move the downloaded binary file to a directory in your system's `PATH`.
3. Rename the binary file to a more convenient name, if desired.

## Build Instructions

If you prefer to build the program from source, follow these instructions:

1. Clone the Github repository to your local machine:
   ```
   git clone https://github.com/vpoluyaktov/abb_ia.git
   ```
2. Ensure you have Go installed on your system.
3. Navigate to the cloned repository directory:
   ```
   cd abb_ia
   ```
4. Build the program using the following command:
   ```
   go build -o bin/
   ```
5. The binary file (`abb_ia`) will be generated in the bin/ directory.

## Reporting Bugs through GitHub Issues

If you encounter any bugs, issues, or unexpected behavior while using **Audiobook Builder**, please follow these steps to report them using GitHub Issues:

1. Visit the [GitHub Issues](https://github.com/vpoluyaktov/abb_ia/issues) page of this repository.

2. Click on the "New Issue" button to create a new issue.

3. In the issue title, provide a concise and descriptive summary of the bug or issue you have encountered.

4. In the issue description, provide detailed information about the bug or issue. Be sure to include the following:
   - Steps to reproduce the bug or issue.
   - Expected behavior.
   - Actual behavior observed.
   - Screenshots or error messages, if applicable.

5. If possible, try to isolate the issue and provide a minimal, reproducible example. This will greatly help us in identifying and resolving the problem.

6. Assign an appropriate label to the issue. You can choose from existing labels or create a new one if necessary.

7. If you have any additional information or suggestions related to the bug or issue, feel free to add them in the issue description.

8. Once you have filled in all the necessary information, click on the "Submit new issue" button to create the bug report.

## Keeping Track of Bug Reports

To keep track of reported bugs and their resolutions, we will regularly review and update the Issues section. You can view the status of your reported bugs, as well as any updates or conversations related to them.

## Contributing to Bug Fixes

If you are interested in contributing to the bug fixes for **Audiobook Builder**, please follow the following steps:

1. Fork the **abb_ia** repository.

2. Create a new branch with a descriptive name for your bug fix.

3. Make the necessary changes to fix the bug.

4. Commit and push your changes to your forked repository.

5. Create a pull request to merge your bug fix branch into the main **abb_ia** repository.

6. Your contribution will be reviewed by the project maintainer, and any necessary feedback or changes will be communicated via the pull request.

## Author's Disclaimer

Since the copyrights for the majority of old-time radio shows have expired and many of them are now in the Public Domain, you have the ability to freely download and listen to them. However, it's important to note that there is also copyrighted content available on the Internet Archive site. It is crucial that you respect the legal rights of others and refrain from engaging in any unlawful activities. This script serves as a tool solely for the purpose of creating an audiobook. The author holds no responsibility for how you choose to utilize this tool. It is your responsibility to adhere to the terms outlined in the copyright license of each item.

## Todo
- The Text to Speech (**TTS**) version of Audiobook Builder is coming soon. It will allow you to create audiobooks from .epub, fb2, and other formats of electronic books.


