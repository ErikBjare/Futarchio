angular.module('WeArePeopleApp', ["ngResource", "ngRoute"])

.factory('msgStack', function() {
    var msgs = [];
    return function(obj) {
        if(msg === "") return msgs;
        msgs.push({type: "info", msg: msg});
        console.log(obj);
    };
})

.factory('User', function($log, $resource) {
    var User = $resource('/api/0/:resource/:key/:val', {});
    return function(key, val) {
        var user = User.get(
           {"resource": "users",
            "key": key,
            "val": val},
            function(u) {
                $log.info(u.data[0]);
            }, function(u) {
                $log.error("Error");
            });
        return user;
    };
})

.controller('MainController', function($scope, $route, $location) {
    $scope.pages = [{title: "Polls", url: "polls"},
                    {title: "Predictions", url: "predictions"},
                    {title: "Admin", url: "admin"}];
    $scope.location = $location;
})

.controller('PageController', function($scope, $log) {
    $scope.log = function(obj) {
    };
})

.controller('AdminController', function($scope, $resource, msgStack, User) {
    $scope.resources = ["users", "idklol"];
    $scope.keys = ["email", "id"];
    $scope.stdin = "";

    $scope.log = function(obj) {
        $scope.stdout = msgStack($scope, obj);
        $scope.stdin = "";
    };

    $scope.findUser = function() {
        var user = User($scope.lookupKey, $scope.lookupValie);
        console.log(user);
    };
})

.controller('PollsController', function($scope, $resource, $log) {
    $scope.polls = [{"user": "Someone", "text": "This is a poll"},
                    {"user": "Someone else", "text": "This is another poll"}];
})

.controller('PollController', function($scope, $resource, $log) {
    $scope.createPoll = function() {
        $log.warn("Not implemented");
    };
})

.controller('ProfileController', function($scope, $routeParams, User) {
    $scope.user = {};
    User("id", $routeParams.userId).$promise.then(function(payload) {
        console.log(payload);
        $scope.user = payload.data[0];
    });
})

.config(function($routeProvider, $locationProvider) {
  $routeProvider
   .when('/polls', {
    templateUrl: 'polls.html',
    controller: 'PollsController',
  })
   .when('/profile/:userId', {
    templateUrl: 'profile.html',
    controller: 'ProfileController',
  })
   .when('/admin', {
    templateUrl: 'admin.html',
    controller: 'AdminController',
  });
});
