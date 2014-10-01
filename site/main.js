angular.module('WeArePeopleApp', ["ngResource", "ngRoute"])
.controller('MainController', function($scope, $route, $location) {
    $scope.pages = [{title: "Polls", url: "polls"},
                    {title: "Predictions", url: "predictions"},
                    {title: "Admin", url: "admin"}];
    $scope.location = $location;
})
.controller('AdminController', function($scope, $resource) {
    $scope.resources = ["users", "idklol"];
    $scope.keys = ["email", "id"];
    $scope.stdout = [];
    $scope.stdin = "";

    $scope.log = function(obj) {
        if(obj !== undefined) {
            msg = obj;
        } else {
            msg = $scope.stdin;
        }
        if(msg === "") return;
        $scope.stdout.push({type: "info", msg: msg});
        console.log($scope.stdin);
        $scope.stdin = "";
    };

    $scope.findUser = function() {
        var User = $resource('/api/0/:resource/:key/:val', {});
        var user = User.get(
           {"resource": "users",
            "key": $scope.lookupKey,
            "val": $scope.lookupValue},
            function(u) {
                $scope.log(u.data[0]);
            }, function(u) {
                $scope.log("Error");
            });
        console.log(user);
    };
})
.controller('PollsController', function($scope, $resource) {
    $scope.polls = [{"user": "Someone", "text": "This is a poll"},
                    {"user": "Someone else", "text": "This is another poll"}];
    $scope.createPoll = function() {
        $scope.log("Not implemented");
    };
})
.config(function($routeProvider, $locationProvider) {
  $routeProvider
   .when('/polls', {
    templateUrl: 'polls.html',
    controller: 'PollsController',
    resolve: {
      // I will cause a 1 second delay
      delay: function($q, $timeout) {
        var delay = $q.defer();
        $timeout(delay.resolve, 1000);
        return delay.promise;
      }
    }
  })
   .when('/admin', {
    templateUrl: 'admin.html',
    controller: 'AdminController',
    resolve: {
      // I will cause a 1 second delay
      delay: function($q, $timeout) {
        var delay = $q.defer();
        $timeout(delay.resolve, 1000);
        return delay.promise;
      }
    }
  });
});
