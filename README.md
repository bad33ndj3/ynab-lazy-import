# ing-to-ynab

### What does it do

ING to YNAB will go to your download directory and try to find csv exports from ING Bank
It can find these because they have a filename that matches your IBAN.

It will simulate a csv import to YNAB by reading these files and uploading them with the `csv Imported` tag.
This way the upload will behave just like a CSV Upload to YNAB.

This approach is used because ING doesn't offer a consumer based API.

### usage

Fill in the env file

    $ cp .env.example .env
    
Run the tool

    $ make run
    
### Supported banks

Currently, this only supports the Dutch ING Bank exports.
