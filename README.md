# Excel Transactions

This is a Go program that reads an Excel file containing transactions and calculates the Adjusted Cost Base and Capital Gains for each symbol.

## Usage

1. Install Go.
2. Clone this repository.
3. Install the required dependencies by running `go get`.
4. Create an Excel file containing transactions. The file should have the following columns:

   - Symbol
   - Date Time (in format "YYYY-DD-MM, hh:mm:ss")
   - Share Quantity
   - Price per share
   - Proceeds / Amounts
   - Commission fee

5. Run the program by running `go run main.go <filename>`, where `<filename>` is the name of the Excel file.
6. The program will output the Adjusted Cost Base and Capital Gains for each symbol in JSON format.

## Example

```
$ go run main.go transactions.xlsx
[
    {
        "Symbol": "AAPL",
        "Qty": 100,
        "ACB": "150.00",
        "Capital Gains": "0.00"
    },
    {
        "Symbol": "GOOG",
        "Qty": 50,
        "ACB": "500.00",
        "Capital Gains": "0.00"
    }
]
```

## Dependencies

## License

This program is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
```

You can modify this template to suit your needs.

I hope this helps! Let me know if you have any questions.
