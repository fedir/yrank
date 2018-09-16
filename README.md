# yrank = Youtube Rank

Package which helps You to priorize a Youtube's playlist items for watching.

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

## Installation

    go get -u yrank

## Usage

IMPORTANT: You should have a Youtube API key to use the application. You could crate it at [Google Developers Console](https://console.developers.google.com/).

To launch the application You should use

    yrank PLAYLIST-ID

Or :

    yrank PLAYLIST-URL

## Advanced configuration

Interestingness formula coefficients :

* Popularity effect : Ordering by number of views, multiplied by likes / dislikes proportion
* Feedback : Ordering by number of comments

    Ci = (Nviews + Nlikes + Ndislikes - Nuniqcomments) / Nviews
