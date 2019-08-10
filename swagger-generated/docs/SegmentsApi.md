# \SegmentsApi

All URIs are relative to *https://www.strava.com/api/v3*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ExploreSegments**](SegmentsApi.md#ExploreSegments) | **Get** /segments/explore | Explore segments
[**GetLeaderboardBySegmentId**](SegmentsApi.md#GetLeaderboardBySegmentId) | **Get** /segments/{id}/leaderboard | Get Segment Leaderboard
[**GetLoggedInAthleteStarredSegments**](SegmentsApi.md#GetLoggedInAthleteStarredSegments) | **Get** /segments/starred | List Starred Segments
[**GetSegmentById**](SegmentsApi.md#GetSegmentById) | **Get** /segments/{id} | Get Segment
[**StarSegment**](SegmentsApi.md#StarSegment) | **Put** /segments/{id}/starred | Star Segment


# **ExploreSegments**
> ExplorerResponse ExploreSegments(ctx, bounds, optional)
Explore segments

Returns the top 10 segments matching a specified query.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **bounds** | [**[]float32**](float32.md)| The latitude and longitude for two points describing a rectangular boundary for the search: [southwest corner latitutde, southwest corner longitude, northeast corner latitude, northeast corner longitude] | 
 **optional** | ***ExploreSegmentsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ExploreSegmentsOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **activityType** | **optional.String**| Desired activity type. | 
 **minCat** | **optional.Int32**| The minimum climbing category. | 
 **maxCat** | **optional.Int32**| The maximum climbing category. | 

### Return type

[**ExplorerResponse**](ExplorerResponse.md)

### Authorization

[strava_oauth](../README.md#strava_oauth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetLeaderboardBySegmentId**
> SegmentLeaderboard GetLeaderboardBySegmentId(ctx, id, optional)
Get Segment Leaderboard

Returns the specified segment leaderboard.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **int64**| The identifier of the segment leaderboard. | 
 **optional** | ***GetLeaderboardBySegmentIdOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetLeaderboardBySegmentIdOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **gender** | **optional.String**| Filter by gender. | 
 **ageGroup** | **optional.String**| Summit Feature. Filter by age group. | 
 **weightClass** | **optional.String**| Summit Feature. Filter by weight class. | 
 **following** | **optional.Bool**| Filter by friends of the authenticated athlete. | 
 **clubId** | **optional.Int64**| Filter by club. | 
 **dateRange** | **optional.String**| Filter by date range. | 
 **contextEntries** | **optional.Int32**|  | 
 **page** | **optional.Int32**| Page number. | 
 **perPage** | **optional.Int32**| Number of items per page. Defaults to 30. | [default to 30]

### Return type

[**SegmentLeaderboard**](SegmentLeaderboard.md)

### Authorization

[strava_oauth](../README.md#strava_oauth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetLoggedInAthleteStarredSegments**
> []SummarySegment GetLoggedInAthleteStarredSegments(ctx, optional)
List Starred Segments

List of the authenticated athlete's starred segments. Private segments are filtered out unless requested by a token with read_all scope.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***GetLoggedInAthleteStarredSegmentsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a GetLoggedInAthleteStarredSegmentsOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **optional.Int32**| Page number. | 
 **perPage** | **optional.Int32**| Number of items per page. Defaults to 30. | [default to 30]

### Return type

[**[]SummarySegment**](SummarySegment.md)

### Authorization

[strava_oauth](../README.md#strava_oauth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetSegmentById**
> DetailedSegment GetSegmentById(ctx, id)
Get Segment

Returns the specified segment. read_all scope required in order to retrieve athlete-specific segment information, or to retrieve private segments.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **int64**| The identifier of the segment. | 

### Return type

[**DetailedSegment**](DetailedSegment.md)

### Authorization

[strava_oauth](../README.md#strava_oauth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **StarSegment**
> DetailedSegment StarSegment(ctx, id, starred)
Star Segment

Stars/Unstars the given segment for the authenticated athlete. Requires profile:write scope.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **int64**| The identifier of the segment to star. | 
  **starred** | **bool**| If true, star the segment; if false, unstar the segment. | [default to false]

### Return type

[**DetailedSegment**](DetailedSegment.md)

### Authorization

[strava_oauth](../README.md#strava_oauth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

