# \RunningRacesApi

All URIs are relative to *https://www.strava.com/api/v3*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetRunningRaceById**](RunningRacesApi.md#GetRunningRaceById) | **Get** /running_races/{id} | Get Running Race
[**GetRunningRaces**](RunningRacesApi.md#GetRunningRaces) | **Get** /running_races | List Running Races


# **GetRunningRaceById**
> RunningRace GetRunningRaceById(ctx, id)
Get Running Race

Returns a running race for a given identifier.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **int32**| The identifier of the running race. | 

### Return type

[**RunningRace**](RunningRace.md)

### Authorization

[strava_oauth](../README.md#strava_oauth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetRunningRaces**
> []RunningRace GetRunningRaces(ctx, optional)
List Running Races

Returns a list running races based on a set of search criteria.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GetRunningRacesOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetRunningRacesOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **year** | **optional.Int32**| Filters the list by a given year. | 

### Return type

[**[]RunningRace**](RunningRace.md)

### Authorization

[strava_oauth](../README.md#strava_oauth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

