# ynab-lazy-import

### What does it do

ynab-lazy-import will go to your download directory and try to find csv exports from ING Bank
ing-to-ynab will go to your download directory and try to find csv exports from ING Bank
It can find these because they have a filename that matches your IBAN.

It will simulate a csv import to YNAB by reading these files and uploading them with the `csv Imported` tag.
This way the upload will behave just like a CSV Upload to YNAB.

This approach is used because ING doesn't offer a consumer based API.

### Usage
Install package

    $ go get -d github.com/bad33ndj3/ynab-lazy-import && go install github.com/bad33ndj3/ynab-lazy-import

Initialize the config file and fill in the missing fields

    $ ynab-lazy-import init -t <access token>
    
Scrape the Downloads folder and upload matching csv's

    $ ynab-lazy-import api
    
### Supported banks

Currently, this only supports the Dutch ING Bank exports.

### Supported platforms
tested on:
- mac
- ubuntu