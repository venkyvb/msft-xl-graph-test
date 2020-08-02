## msft-xl-graph-test

This test program uses MSFT Graph API for Excel to populate a set of input cells, which drives some calcualtion in the XL sheet and read the results from a set of output cells. The cells are hard-coded in the program (as this is intended for just test purposes).

In order to run the test a config file either in `yaml` or `json` format should be provided as input. THe structure of the `yaml` file is as shown below.

To run the test

```
msft-xl-graph-test run --config <path_to_config_file>
```
If the `--config` option is not specified the test looks for the config file under the HOME directory with the file name `.msft-xl-graph-test.yaml`.

```
accessToken: <access_token_from_graph_explorer>
workbookItemID: <workbookItemID_for_the_test_workbook>
noOfIterations: 5
inputParams:
  - memCnt: 230
    recCnt: 79
    curr: USD
  - memCnt: 1230
    recCnt: 79
    curr: CAD
  - memCnt: 12300
    recCnt: 260
    curr: GBP
  - memCnt: 36900
    recCnt: 749
    curr: EUR
  - memCnt: 2500
    recCnt: 441
    curr: AUD    
  - memCnt: 2500
    recCnt: 441
    curr: USD      
```    
