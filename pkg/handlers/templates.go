package handlers

// I know this is dummy but I didn't want to fight with GAE to have
// html templates in custom place.

var IndexHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <title>Title</title>
        <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css"
              integrity="sha384-MCw98/SFnGE8fJT3GXwEOngsV7Zt27NXFoaoApmYm81iuXoPkFOJwJ8ERdknLPMO" crossorigin="anonymous">
    </head>
    <style>
        .center {
            display: block;
            margin-left: auto;
            margin-right: auto;
            width: 50%;
        }
    </style>
    <body>
        <div class="container">
            <div class="row" >
                <div class="col-12 align-self-center">
                        <h1>
                           <p class="text-center font-weight-light">Plot your activities with plotted!</p>
                        </h1>
                        ​<picture class="text-center">
                            <a class="center" href="{{.Auth}}">
                                <img src="static/btn_strava_connectwith_orange.svg" class="shadow p-3 mb-5 bg-white rounded"/>
                            </a>
                        </picture>
                </div>
            </div>
        </div>
    </body>
</html>
`

var MapHTML = `
<html>
  <head>
    <link rel="stylesheet" href="https://unpkg.com/leaflet@1.3.1/dist/leaflet.css"
    integrity="sha512-Rksm5RenBEKSKFjgI3a41vrjkw4EVPlJ3+OiI65vTjIdo9brlAacEuKOiQ5OFh7cOI1bkDwLqdLw3Zg0cRJAAQ=="
    crossorigin=""/>
      <!-- Make sure you put this AFTER Leaflet's CSS -->
    <script src="https://unpkg.com/leaflet@1.3.1/dist/leaflet.js"
    integrity="sha512-/Nsx9X4HebavoBvEBuyp3I7od5tA0UzAxs+j83KgC8PU0kgB4XiK4Lfe4y4cgBtaRJQEIFCW+oC506aPT2L1zw=="
    crossorigin=""></script>
    <script type="text/javascript" src="https://rawgit.com/jieter/Leaflet.encoded/master/Polyline.encoded.js"></script>
    <style>
      #mapid { height: 100%; }
    </style>
  </head>
  <body>
    <div id="mapid"></div>
    <script>
      function rand(min, max) {
        return parseInt(Math.random() * (max-min+1), 10) + min;
      }
      function getRandomColor() {
        var h = rand(200, 220);
        var s = rand(70, 100);
        var l = rand(20, 50);
        return 'hsl(' + h + ',' + s + '%,' + l + '%)';
      }

        var mymap = L.map('mapid').setView([52.380000, 16.920000] , 12);
        L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token={{.MapboxToken}}', {
            maxZoom: 18,
            id: 'mapbox.streets',
        }).addTo(mymap);
        var encodedRoutes = [
          {{range .EncodedRoutes}}
              {{.}},
          {{end}}
        ];
        for (let route of encodedRoutes) {
          L.polyline(
            route,
            {
              color: getRandomColor(),
              weight: 3,
              opacity: .9,
              lineJoin: 'round',
            }
          ).addTo(mymap);
        }
      </script>
  </body>
</html>
`
