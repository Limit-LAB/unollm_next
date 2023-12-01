# \KeyStoreAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**KeyStoreRemoveMapToPost**](KeyStoreAPI.md#KeyStoreRemoveMapToPost) | **Post** /keyStore/removeMapTo | 删除现有映射
[**KeyStoreTestTransformerPost**](KeyStoreAPI.md#KeyStoreTestTransformerPost) | **Post** /keyStore/testTransformer | 测试APIKey映射
[**KeyStoreUserDefinedKeysGet**](KeyStoreAPI.md#KeyStoreUserDefinedKeysGet) | **Get** /keyStore/userDefinedKeys | 获取用户定义API Key表列



## KeyStoreRemoveMapToPost

> map[string]interface{} KeyStoreRemoveMapToPost(ctx).MapTo(mapTo).Execute()

删除现有映射



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
    mapTo := *openapiclient.NewMapTo([]int32{int32(123)}, int32(123)) // MapTo |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.KeyStoreAPI.KeyStoreRemoveMapToPost(context.Background()).MapTo(mapTo).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `KeyStoreAPI.KeyStoreRemoveMapToPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `KeyStoreRemoveMapToPost`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `KeyStoreAPI.KeyStoreRemoveMapToPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiKeyStoreRemoveMapToPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mapTo** | [**MapTo**](MapTo.md) |  | 

### Return type

**map[string]interface{}**

### Authorization

[apikey-header-Authorization](../README.md#apikey-header-Authorization)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## KeyStoreTestTransformerPost

> KeyStoreTestTransformerPost200Response KeyStoreTestTransformerPost(ctx).KeyStoreTestTransformerPostRequest(keyStoreTestTransformerPostRequest).Execute()

测试APIKey映射



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
    keyStoreTestTransformerPostRequest := *openapiclient.NewKeyStoreTestTransformerPostRequest("Key_example", "Provider_example") // KeyStoreTestTransformerPostRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.KeyStoreAPI.KeyStoreTestTransformerPost(context.Background()).KeyStoreTestTransformerPostRequest(keyStoreTestTransformerPostRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `KeyStoreAPI.KeyStoreTestTransformerPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `KeyStoreTestTransformerPost`: KeyStoreTestTransformerPost200Response
    fmt.Fprintf(os.Stdout, "Response from `KeyStoreAPI.KeyStoreTestTransformerPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiKeyStoreTestTransformerPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **keyStoreTestTransformerPostRequest** | [**KeyStoreTestTransformerPostRequest**](KeyStoreTestTransformerPostRequest.md) |  | 

### Return type

[**KeyStoreTestTransformerPost200Response**](KeyStoreTestTransformerPost200Response.md)

### Authorization

[apikey-header-Authorization](../README.md#apikey-header-Authorization)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## KeyStoreUserDefinedKeysGet

> []KeyStoreUserDefinedKeysGet200ResponseInner KeyStoreUserDefinedKeysGet(ctx).Execute()

获取用户定义API Key表列



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
    resp, r, err := apiClient.KeyStoreAPI.KeyStoreUserDefinedKeysGet(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `KeyStoreAPI.KeyStoreUserDefinedKeysGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `KeyStoreUserDefinedKeysGet`: []KeyStoreUserDefinedKeysGet200ResponseInner
    fmt.Fprintf(os.Stdout, "Response from `KeyStoreAPI.KeyStoreUserDefinedKeysGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiKeyStoreUserDefinedKeysGetRequest struct via the builder pattern


### Return type

[**[]KeyStoreUserDefinedKeysGet200ResponseInner**](KeyStoreUserDefinedKeysGet200ResponseInner.md)

### Authorization

[apikey-header-Authorization](../README.md#apikey-header-Authorization)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

