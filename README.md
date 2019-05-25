# test_task_15
Some test task

See description of task in:
```
 ./task_desc
```

Run:
```
go run cmd/loaner/*.go
```

See result (after run):
```
out/assignments.csv
out/yields.csv
```

See My Result:
```
my_out/assignments.csv
my_out/yields.csv
```

Answers on Questions:
1. I spend ~ 8 hours; Understanding all requirements.
1. I will extend Model Covenant with adding new necessary fields and also extend helper model FacilCovens and BankCovens. And finally will change validateFacility function for handle new rules.
1. We can pass new Facilities through facilities update buffer . One way is We can before start handle new loan we can update facilities from this buffer. Or we can after finish handle loan check is there updates if yes repeat loan handling after apply updates. Also we can break loan handling if come update and start process again after updating. First most easy but can work with not actual data, second and third solution always will work with actual data.
1. REST API example:
    ```
    Request Loan (send request to get loan, put queue/stream)
    POST /request_loan
    BODY: {...}
    Response: {"request_id"}

    Request Status of Loan Request
    GET /request_loan/{request_id}/status
    Response: {"status": "Processing"} or {"status": "OK", "facility": {...}}
    ```
1. ?
1. Runtime Complexity:
Time Complexity in simple explanation (no include parse inputs CSV and build models) is:
O(L x F x C(F) x C(B(F))) where L - len of Loan items in stream, F - len of Facilities, C(F) - count of covenants by facility (used additional map for store covenants by facilities), C(B(F)) - covenants of facility bank. Because we use C(F) and C(B(F)) we no need iterate full Covenant list. Also because F sorted but rate, in most case we don't look full list of Facilities by each Loan, so it is also decrease our complexity.
Memory complexity: O(L x F x C x CFM x CBFM x FCap x A x Y x FY), L - loan, F - facilities, C - covenants, CFM - covenants by facilities map, CBFM - covenants by facilities banks map, FCap - facilities actual capacity map, A - assignments out, Y - yield out, FY - facilities yield map.
I think there is possible reduce complexity by optimisation working with facilities, for example not to look at facilities that are no longer relevant (but it will increase memory complexity).