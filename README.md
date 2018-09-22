# yrank = Youtube Rank

Package which helps You to priorize a Youtube's playlist items for watching.

It could be quite helpful, when You would like to choose the most interesting videos of an IT conference.

## Installation

golang 1.10+ should be installed. Go environment [is very simple to install](https://golang.org/doc/install).

To install the package :

    go get -u github.com/fedir/yrank

## Usage

IMPORTANT: You should have a Youtube API key to use the application. You could crуate it at [Google Developers Console](https://console.developers.google.com/).

After that, please copy ```config.example.toml``` to ```config.toml``` in the application folder and define Your Youtube API key there.

To launch the application You should just precise the ID of the playlist via CLI (this ID could be found in the URL of the playlist, it's the "?playlistId=" variable value).

    yrank -p PLAYLIST-ID

To output ranking in markdown:

    yrank -p PLAYLIST-ID -o markdown

### Example

    ./yrank -p PL2ntRZ1ySWBdatAqf-2_125H4sGzaWngM

### Results samples

* Ranking of GopherCon 2018 videos https://gist.github.com/fedir/98f6a2ed65e7462a101198dc6f3d5185
* Ranking of GopherCon UK 2018 videos https://gist.github.com/fedir/6a93e91fa414df6484ba04589ed3269a
