# Publit Production API SDK for GO

ProductionAPIGoSDK is the official Go SDK for the Publit production API.

The SDK aims to help implementation against the Publit production API by supplying Go libraries for handling 
calls and structs against the API.

The ProductionAPIGoSDK contains functionality for print order handling and automation.

## Installing

Retrieve the SDK by running:

`$ go get github.com/publitsweden/ProductionAPIGoSDK`

### Dependencies

The SDK has dependencies to the APIUtilityGoSDK which contains common heplers for the Publit APIs.

## Usage

See the Godocs: https://golang.org/pkg/github.com/publitsweden/ProductionAPIGoSDK for more information about implementation, examples and usage.

**Simple example**
```Go
// Create APIClient (note incomplete for breivity).
c := production.APIClient{ ... }

// Call "Show" with the client and a print order ID to retrieve that order from the Publit API.
p, _ : printorder.Show(c, 1)
// Do something with "p" (the returned print order).
```

### production.APIClient
The packages in the ProductionAPIGoSDK have various method for getting, updating, storing and deleting data in Publit through the Production API.

Each of these methods in the packages have interfaces dictating how the call against the API should be made.
The production.APIClient fulfils these interfaces and helps with authorisation against the Publit API.

Therefore most usages of the SDK should begin with creating a production.APIClient.

```Go
c := production.APIClient{}
```

The APIClient has a Client attribute which takes a *common.Client, which performs the actual calls against the API.

It also takes the base URL as argument. This can be used for changing between Publit production and sandbox environments.

The below snippet shows a more in depth example:

```Go
c := production.APIClient{
        Client: *client.Client{},//See the APIUtilityGoSDK for more information about client.Client
        BaseUrl: "https://url.to.publit",
}
```

## Examples
The examples under this section serves only as illustrative examples on how to use the ProductionAPIGoSDK.

**Fetching print orders**

Below is an example on how to fetch print orders after a certain dat filtered by a status. 
The response is also limited to 1.

```Go 
//Create a production.APIClient that will aid in making calls against the API.
c := production.APIClient{
    Client:  client.New(
         func(c *client.Client) {
             c.User = "myusername"
             c.Password = "mypassword"
         },
    ),
    BaseUrl: "https://url.to.publit",
}

// Check if service is up.
ok := c.StatusCheck()

if !ok {
    fmt.Println("Status check not ok")
} else {
    log.Fatal("Status check ok")
}

// Create a filter on the attribute "created_at".
// The filter will only fetch results that has been created after or equal to the 1st ofh january 2017.
createdAtFilter := common.AttrQuery{
    Name: printorder.CREATED_AT,
    Value: "2017-01-01 00:00:00",
    Args: common.AttrArgs{
        Operator: common.OPERATOR_GREATER_EQUAL,
        Combinator: common.COMBINATOR_AND,
    },
}

// Create a filter on the attribute "status".
// The filter will only fetch results that has the status "exported". The state constant resides in the printorderstatus package.
statusFilter := common.AttrQuery{
    Name: "status",
    Value: printorderstatus.STATE_EXPORTED.AsString(),
}

// Create an Index request on the printorder package.
// It takes a ProductionAPIClient c (or any client that fulfils the interface states for the Index method in the printorders package).
// It also takes a list of variadic functions, which the common library has helpers for compiling more easily.
// The index method returns a printorder.IndexResponse.
pos, err := printorder.Index(
    c,
    common.QueryLimit(1, 0),// Add the limit.
    common.QueryWith(printorder.WITH_STATUSES),// Add a with parameter (this loads the resource with any related printorderstatuses). 
    common.QueryAttr(createdAtFilter),// Add the createdAt filter.
    common.QueryAttr(statusFilter),// Add the status filter.
    common.QueryOrderBy([]string{"id"}, common.ORDER_DIR_DESC),// Order results by id in descending direction.
)

// Handle error.
if err != nil {
    log.Fatal(err.Error())
}

// Range through the "data" attribute in the IndexResponse to output received print orders.
for _, v := range pos.Data {
    fmt.Println("Found print orders:")
    fmt.Printf("%+v\n", v)
}
```

**Setting print order status**

Below is an example on how to set a status for an order. Note that this is only an illustrative example.
A real program would most likely derive the order that the status should be set to from some ingestion method (like the one above).

```Go
//Create a production.APIClient that will aid in making calls against the API.
c := production.APIClient{
    Client:  client.New(
         func(c *client.Client) {
             c.User = "myusername"
             c.Password = "mypassword"
         },
    ),
    BaseUrl: "https://url.to.publit",
}

// Check if service is up.
ok := c.StatusCheck()

if !ok {
    fmt.Println("Status check not ok")
} else {
    log.Fatal("Status check is not ok")
}

orderId := 12345
// Create a new status with a "state", an order id and an optional message.
s := printorderstatus.New(printorderstatus.STATE_ACCEPTED, orderId, "")
    
// Run store. Store will also update "s" with data retrieved from the API.
err := s.Store(c)

if err != nil {
        log.Fatal(err.Error())
}

// Print s
fmt.Printf("%+v",s)

```

**Downloading files**

Below is an example on how to download print order files from the API.

```Go
//Create a production.APIClient that will aid in making calls against the API.
c := production.APIClient{
    Client:  client.New(
         func(c *client.Client) {
             c.User = "myusername"
             c.Password = "mypassword"
         },
    ),
    BaseUrl: "https://url.to.publit",
}

// Create a list of files. Can also be retrieved from the API.
fl := file.FileList{
    &file.File{ID: 1234},
    &file.File{ID: 5678},
}

// Output destination folder.
output := "/path/to/folder/"

// Download files.
errList, err := fl.DownloadFiles(c, output)

// Check if method returned error.
if err != nil {
        log.Fatal(err.Error())
}

// Check if any errors occured when downloading each item in the FileList.
// The errList is a map[int]error where the map is indexed by File.ID.
for k, v := range errList {
    if v != nil {
        fmt.Sprintf("Error for fileID: %v", k)
        fmt.Println("Error: ", v.Error())
    }
}

```