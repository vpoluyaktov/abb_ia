# Audiobook Builder (Internet Archive version)

## Description

The Internet Archive site offers a vast collection of free "old-time radio" shows, audiobooks, and lectures that can be downloaded. While you can listen to them on your web browser, it can be inconvenient especially if you want to listen to them on your mobile device. Typically, these shows are provided as individual .mp3 files, requiring you to download them all, create a playlist, and remember your last listened file and position.

To make this process easier, I developed Audiobook Builder. With this app, all you need is the name of a show or book, or a direct link on archive.org. It will download the .mp3 files for the book, re-encode them with the same bit rate, generate a list of chapters (which can be edited during the process), and ultimately create an audiobook in .m4b format.


![Audiobook Builder in action](https://github.com/vpoluyaktov/abb_ia/blob/master/assets/abb_ia_v2.gif)

## Features
- TUI interface. It allows you to run this application either on your own computer or on a remote server using ssh with tmux, screen, or byobu. This can be helpful when creating an audiobook that takes a long time.
- Download a set of single .mp3 files from [archive.org](https://archive.org)
- Create an audiobook in .m4b format
- Re-encode mp3 files to the same bit rate, if necessary.
- Modify audiobook metadata obtained from [archive.org](https://archive.org), including book title, author, series, genre, and art cover
- Copy created audiobook to specified folder located on the same server using [audiobookshelf compatible directory structure](https://www.audiobookshelf.org/docs/#book-directory-structure). This can be helpful when you run `abb_ia` on the same server where the [Audiobookshelf server](https://www.audiobookshelf.org) is hosted, or when you mount the Audiobookshelf library folder via NFS.
- Upload your created audiobook to a personal [Audiobookshelf server](https://www.audiobookshelf.org) so that you can easily listen to it on your favorite device.

## Integrations

Audiobook Builder seamlessly integrates with **Audiobookshelf server** (https://www.audiobookshelf.org). This integration allows you to upload the created audiobooks directly to the Audiobook Shelf server for convenient listening.

## Installation Instructions

To use Audiobook Builder, you need to have the following command line utilities installed:

- **ffmpeg** (used for audio manipulation)
- **ffprobe** (used for retrieving audio metadata)

See ffmpeg website for more details: https://ffmpeg.org/

The easiest way to install these utilities on a Linux computer is by running the following command:

```bash
sudo apt install ffmpeg
```

For MacOS, you can use the [Homebrew](https://brew.sh/) utility with the command:

```bash
brew install ffmpeg
```

If you are using Windows, you can find instructions for installation on the **ffmpeg** website: [https://ffmpeg.org/download.html](https://ffmpeg.org/download.html)

Make sure these utilities are properly installed and available in your system's PATH before running `abb_ia`.

To install Audiobook Builder (`abb_ia`) on your system, follow these steps:

1. Download the ready-to-run binary file for your target platform from the [Github Releases page](https://github.com/vpoluyaktov/abb_ia/releases).
2. Move the downloaded binary file to a directory in your system's `PATH`.
3. Open a terminal and navigate to the directory where the binary file is located.
4. Run the binary file by executing the command `./abb_ia`. The TUI interface will appear.
5. Follow the instructions on the application interface to do a search, create an audiobook, and upload it to the [Audiobookshelf server](https://www.audiobookshelf.org) if necessary. <br/>
   You can try searching for:

- **Old Time Radio Researchers Group: Single Episodes**
- **BBC Radio 4: Radio Drama** (make sure to check if the show is copyrighted).
- **CBS Radio: Radio Mystery Theater**
- **Boxcars711:***
- **Greg Lauer:***
- **Relic Radio:***
- **Radio Memories Network:***

6. Enjoy listening to an audiobook on your favorite device.



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


## Known Issues
- On some systems, you may not see a cursor in the input fields. This is because the color of the cursor in your terminal application settings is the same as the background color of the input field. To solve this issue, you can adjust the settings of your terminal program. You can either change the cursor color (so that it is different from the color of the input fields) or make the cursor blink.
- On some Windows systems, the Audiobook builder has a problem launching ffmpeg (there is an issue with the input file path). If you encounter this problem, please enable DEBUG mode in the application settings, replicate the error, and file a GitHub Bug report by attaching the application log file.
- Sometimes, when the application crashes, the terminal window may be filled with random characters and you won't see a cursor anymore. This happens because of a problem with the TUI framework that was used to create the application. To solve this issue, you can either reopen the terminal window or try running the `reset` command.
- Downloading audiobooks using `abb_id` is incredibly easy and fast. This means that over time, you might collect hundreds of audiobooks with thousands of hours of content. However, this can be a problem because it would be practically impossible to listen to all of them in your lifetime. Congratulations!!! You are a data hoarder now! :-)

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

Since the copyrights for the majority of old-time radio shows have expired and many of them are now in the Public Domain, you have the ability to freely download and listen to them. However, it's important to note that there is also copyrighted content available on the Internet Archive site. It is crucial that you respect the legal rights of others and refrain from engaging in any unlawful activities. This application serves as a tool solely for the purpose of creating an audiobook. The author holds no responsibility for how you choose to utilize this tool. It is your responsibility to adhere to the terms outlined in the copyright license of each item.

## Todo
- The Text to Speech (**TTS**) version of Audiobook Builder is coming soon. It will allow you to create audiobooks from .epub, fb2, and other formats of electronic books.

## Join me on Discord
https://discord.gg/ntYyJ7xfzX

