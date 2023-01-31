# PROPER CHALLENGE - SCRAPER

This project aims to fulfill the requirements layed out by the proper team for their take home challenge. The result is a web scraper that scrapes images from the website icanhas.cheezburger.com.

## RUNNING THE PROGRAM

The program takes two optional arguments. The first one, `-amount`, determines the amount of images to be scraped. Each page of the website contains 16 images so every 16 images a new page will be fetched. The second argument, `-threads`, determines the amount of threads the program will use when downloading the images. The default value for the `-amount` argument is 10 while the default value for the `-threads` argument is 5. Amount must be a number greater than zero while threads must be a number between 1 and 5.

An example run would be:

```
> go run cmd/api/main.go -amount 32 -threads 4
```

In order to run all unit tests one can run the command `go test ./...` in the root of the project


## PERFORMANCE

In the root of the project there is a small script, `performance.sh`, that runs the program N number of times for every possible argument of the `-threads` argument (1 through 5). The first argument for the script determines the number of iterations for each thread while the second argument determines the number of images that will be fetched each time the program is called. This script can be used to measure the performance of the program. an example run would be:

```
> ./performance.sh 10 128
1 thread(s): 12.150
2 thread(s): 11.340
3 thread(s): 11.387
4 thread(s): 10.764
5 thread(s): 10.956
```

In general the only noticeable difference seems to be between 1 and >1 threads. With runs with more than 1 thread being approximately a second faster.

## DESIGN

The main logic is in the requesting.go service. The service has 2 exposed methods: `GetImageURL` and `GetImageURLs`. The first one is the callback passed to the colly collector while the second one contains the main logic of the service. the service struct contains a queue where each element is an `action` interface. There are two implementations of this interface `scrapeAction` and `imageAction`. The first implementation calls the `Visit` method of the colly collector while the second implementation downloads an image.

The `GetImageURL` function creates the "results" directory where the images will be stored and proceeds to call each action. When a service is created via the `NewService` function it is created with a single action, a scrape action to page zero of the base url. when a scrape action is called all of the required image actions will be added to the queue and, if there weren't enough images in that page, a new scrape action will be added for the next page.

All of the functions from third packages are abstracted away to interfaces so that the package can be tested. This is the case with the `log` package functions in the logging service, with the `os` package functions in the file service and in the image repository.
