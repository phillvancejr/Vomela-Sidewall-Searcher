# Vomela Sidewall Searcher

A tool created for work. Search sidewall and vinyl color matches. Built with Golang

![screen-gif](demo.gif)

I have no NDAs for tools like this so I decided to upload these tools to github both version control, self reference and for anyone who might be interested to look at them. 

This tool is built with my fork of [webview](https://github.com/phillvancejr/webview) which contains some additional features. This project features:
   - Golang backend to serve the ui content and interface with the OS
   - Simple ui Design via html, css & js

## Updating sidewall.json
sidewall.json in the sidewall_data folder functions as the database for color matches. The original data is stored in several csv files in the same folder which are used to create sidewall.json with go generate. After modifying any csv file make sure to run go generate to update sidewall.json
