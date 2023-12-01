# \DefaultAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**KeyStoreAddKeyPost**](DefaultAPI.md#KeyStoreAddKeyPost) | **Post** /keyStore/addKey | /addKey: 新增API KEY
[**KeyStoreKeysGet**](DefaultAPI.md#KeyStoreKeysGet) | **Get** /keyStore/keys | /keys: 获取全部原始Key
[**KeyStoreMapToPost**](DefaultAPI.md#KeyStoreMapToPost) | **Post** /keyStore/mapTo | /mapTo: 映射APIKey到用户定义Key
[**KeyStoreNewApiPost**](DefaultAPI.md#KeyStoreNewApiPost) | **Post** /keyStore/newApi | /newApi: 新建用户定义的API Key



## KeyStoreAddKeyPost

> KeyStoreAddKeyPost200Response KeyStoreAddKeyPost(ctx).KeyStoreAddKeyPostRequest(keyStoreAddKeyPostRequest).Execute()

/addKey: 新增API KEY



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    keyStoreAddKeyPostRequest := *openapiclient.NewKeyStoreAddKeyPostRequest("Key_example", "Provider_example") // KeyStoreAddKeyPostRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.KeyStoreAddKeyPost(context.Background()).KeyStoreAddKeyPostRequest(keyStoreAddKeyPostRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.KeyStoreAddKeyPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `KeyStoreAddKeyPost`: KeyStoreAddKeyPost200Response
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.KeyStoreAddKeyPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiKeyStoreAddKeyPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **keyStoreAddKeyPostRequest** | [**KeyStoreAddKeyPostRequest**](KeyStoreAddKeyPostRequest.md) |  | 

### Return type

[**KeyStoreAddKeyPost200Response**](KeyStoreAddKeyPost200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## KeyStoreKeysGet

> []KeyStoreKeysGet200ResponseInner KeyStoreKeysGet(ctx).Execute()

/keys: 获取全部原始Key



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.KeyStoreKeysGet(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.KeyStoreKeysGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `KeyStoreKeysGet`: []KeyStoreKeysGet200ResponseInner
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.KeyStoreKeysGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiKeyStoreKeysGetRequest struct via the builder pattern


### Return type

[**[]KeyStoreKeysGet200ResponseInner**](KeyStoreKeysGet200ResponseInner.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## KeyStoreMapToPost

> map[string]interface{} KeyStoreMapToPost(ctx).KeyStoreMapToPostRequest(keyStoreMapToPostRequest).Execute()

/mapTo: 映射APIKey到用户定义Key



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    keyStoreMapToPostRequest := *openapiclient.NewKeyStoreMapToPostRequest([]int32{int32(123)}, int32(123)) // KeyStoreMapToPostRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.KeyStoreMapToPost(context.Background()).KeyStoreMapToPostRequest(keyStoreMapToPostRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.KeyStoreMapToPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `KeyStoreMapToPost`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.KeyStoreMapToPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiKeyStoreMapToPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **keyStoreMapToPostRequest** | [**KeyStoreMapToPostRequest**](KeyStoreMapToPostRequest.md) |  | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## KeyStoreNewApiPost

> map[string]interface{} KeyStoreNewApiPost(ctx).KeyStoreNewApiPostRequest(keyStoreNewApiPostRequest).Execute()

/newApi: 新建用户定义的API Key



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    keyStoreNewApiPostRequest := *openapiclient.NewKeyStoreNewApiPostRequest() // KeyStoreNewApiPostRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.DefaultAPI.KeyStoreNewApiPost(context.Background()).KeyStoreNewApiPostRequest(keyStoreNewApiPostRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.KeyStoreNewApiPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `KeyStoreNewApiPost`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.KeyStoreNewApiPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiKeyStoreNewApiPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **keyStoreNewApiPostRequest** | [**KeyStoreNewApiPostRequest**](KeyStoreNewApiPostRequest.md) |  | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

