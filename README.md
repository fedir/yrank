# yrank = Youtube Rank

Package which helps You to priorize a Youtube's playlist items for watching.

It could be quite helpful, when You would like to choose the most interesting videos of an IT conference.

## Installation

golang 1.10+ should be installed. Go environment [is very simple to install](https://golang.org/doc/install).

To install the package :

    go get -u yrank

## Usage

IMPORTANT: You should have a Youtube API key to use the application. You could crуate it at [Google Developers Console](https://console.developers.google.com/).

After that, please copy ```config.example.toml``` to ```config.toml``` in the application folder and define Your Youtube API key there.

To launch the application You should just precise the ID of the playlist via CLI (this ID could be found in the URL of the playlist, it's the "?playlistId=" variable value).

    yrank -p PLAYLIST-ID

### Example

    ./yrank -p PL2ntRZ1ySWBdatAqf-2_125H4sGzaWngM

### Results samples

* Ranking of GopherCon 2018 videos https://gist.github.com/fedir/98f6a2ed65e7462a101198dc6f3d5185
* Ranking of GopherCon UK 2018 videos https://gist.github.com/fedir/6a93e91fa414df6484ba04589ed3269a

---

# TBD

## Interestingness coefficient formula

    Ci = (Nviews + Nlikes - Ndislikes + Nuniqcomments) / Nviews

### Example of real video statistics

    "viewCount": "57692724",
    "likeCount": "213510",
    "dislikeCount": "44399",
    "commentCount": "6881"

The interestingness coefficient of this video will be :

    Ci = (213510 - 44399 + 6881) / 57692724 = 0.003

### Theoreticale examples of interistingness coefficients

Let's suppose everybody likes the video, nobody dislikes, everybody let's comments, all comments are unique.

    Ci = (57692724 - 0 + 57692724) / 57692724 = 2

Let's suppose nobody likes the video, every watcher dislikes it, everybody let's comments, all comments are unique.

    Ci = (0 - 57692724 + 57692724) / 57692724 = 1

We could notice, what if the video has not at all any likes, the coefficient is still some high. It explains by following idea: even if the video is not ranked, it has lot's of comments, what's mean, it's in some kind interesting.

### Coefficient optimization

Also it's possible in future versions of yrank to make coefficients and formula to be configured through the configuration system, so everybody could choose his own coefficients. At the moment we'll keep simple solution.

## Advanced configuration

Interestingness formula coefficients :

* Popularity effect : Ordering by number of views, multiplied by likes / dislikes proportion
* Feedback : Ordering by number of comments

    Ci = (Nviews + Nlikes + Ndislikes - Nuniqcomments) / Nviews
