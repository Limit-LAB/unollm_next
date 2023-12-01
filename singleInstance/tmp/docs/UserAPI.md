# \UserAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**UserCreatePost**](UserAPI.md#UserCreatePost) | **Post** /user/create | 创建用户
[**UserLoginPost**](UserAPI.md#UserLoginPost) | **Post** /user/login | 登陆
[**UserLogoutGet**](UserAPI.md#UserLogoutGet) | **Get** /user/logout | 登出



## UserCreatePost

> map[string]interface{} UserCreatePost(ctx).UserCreatePostRequest(userCreatePostRequest).Execute()

创建用户



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
    userCreatePostRequest := *openapiclient.NewUserCreatePostRequest("Username_example", "Password_example") // UserCreatePostRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.UserAPI.UserCreatePost(context.Background()).UserCreatePostRequest(userCreatePostRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UserAPI.UserCreatePost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UserCreatePost`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `UserAPI.UserCreatePost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUserCreatePostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userCreatePostRequest** | [**UserCreatePostRequest**](UserCreatePostRequest.md) |  | 

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


## UserLoginPost

> UserLoginPost200Response UserLoginPost(ctx).UserCreatePostRequest(userCreatePostRequest).Execute()

登陆



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
    userCreatePostRequest := *openapiclient.NewUserCreatePostRequest("Username_example", "Password_example") // UserCreatePostRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.UserAPI.UserLoginPost(context.Background()).UserCreatePostRequest(userCreatePostRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UserAPI.UserLoginPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UserLoginPost`: UserLoginPost200Response
    fmt.Fprintf(os.Stdout, "Response from `UserAPI.UserLoginPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUserLoginPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userCreatePostRequest** | [**UserCreatePostRequest**](UserCreatePostRequest.md) |  | 

### Return type

[**UserLoginPost200Response**](UserLoginPost200Response.md)

### Authorization

[apikey-header-Authorization](../README.md#apikey-header-Authorization)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UserLogoutGet

> UserLoginPost200Response UserLogoutGet(ctx).UserCreatePostRequest(userCreatePostRequest).Execute()

登出



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
    userCreatePostRequest := *openapiclient.NewUserCreatePostRequest("Username_example", "Password_example") // UserCreatePostRequest |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.UserAPI.UserLogoutGet(context.Background()).UserCreatePostRequest(userCreatePostRequest).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UserAPI.UserLogoutGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UserLogoutGet`: UserLoginPost200Response
    fmt.Fprintf(os.Stdout, "Response from `UserAPI.UserLogoutGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUserLogoutGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userCreatePostRequest** | [**UserCreatePostRequest**](UserCreatePostRequest.md) |  | 

### Return type

[**UserLoginPost200Response**](UserLoginPost200Response.md)

### Authorization

[apikey-header-Authorization](../README.md#apikey-header-Authorization)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

