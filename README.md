# yrank = Youtube Rank

[![Build Status](https://travis-ci.org/fedir/yrank.svg?branch=master)](https://travis-ci.org/fedir/yrank)
[![codecov](https://codecov.io/gh/fedir/yrank/branch/master/graph/badge.svg)](https://codecov.io/gh/fedir/yrank)

Package which helps You to priorize a Youtube channel or playlist items for watching.

It could be quite helpful, when You would like to choose the most interesting videos of an IT conference.

## Installation

golang 1.10+ should be installed. Go environment [is very simple to install](https://golang.org/doc/install).

To install the package :

    go get -u github.com/fedir/yrank

## Usage

IMPORTANT: You should have a Youtube API key to use the application. You could crуate it at [Google Developers Console](https://console.developers.google.com/).

After that, please copy ```config.example.toml``` to ```config.toml``` in the application folder and define Your Youtube API key there.

### CLI options

    Usage of ./yrank:
    -p string
            Youtube playlist ID
    -c string
            Youtube channel ID
    -s string
        Sorting (default "likes", could be "positive-interest", "total-reaction")
    -o string
            Output format (default "table", could be "markdown")
    -m int
        The maximum number of items that should be returned
    -d bool
            Debug mode for more details during API exchange

### Getting single playlist statistics

To launch the application You should just precise the ID of the playlist via CLI (this ID could be found in the URL of the playlist, it's the "?playlistId=" variable value).

    yrank -p PLAYLIST-ID

To output ranking in markdown:

    yrank -p PLAYLIST-ID -o markdown

### Getting the statistics of a whole user's channel

First of all, You must find the channel ID of the user. It's not always that easy. Sometimes it's in URL of Youtube's user profile. Sometimes, You should look for it in the code of the page.

    yrank -c CHANNEL-ID
    yrank -c CHANNEL-ID -o markdown -s positive-interest

### Examples

    ./yrank -p PL2ntRZ1ySWBdatAqf-2_125H4sGzaWngM
    ./yrank -p PL2ntRZ1ySWBdatAqf-2_125H4sGzaWngM -o markdown -s positive-interest

### Results samples

* Ranking of GopherCon 2018 videos https://gist.github.com/fedir/98f6a2ed65e7462a101198dc6f3d5185
* Ranking of GopherCon UK 2018 videos https://gist.github.com/fedir/6a93e91fa414df6484ba04589ed3269a
* Ranking of Gopher Academy channel videos https://gist.github.com/fedir/c900d0fb59658f9657253f33e38422fe
