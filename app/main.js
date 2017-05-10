var app = angular.module('myApp', []);
app.controller('myCtrl', function($scope, $http) {

    $http.get('http://localhost:3000').
        then(function(response) {
            $scope.people = response.data;
        });
});