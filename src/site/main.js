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

.filter('capitalize', function() {
    return function(input) {
        output = input[0].toUpperCase() + input.substr(1,input.length);
        return output;
    };
})

.controller('MainController', function($scope, $route, $location) {
    $scope.links_left = [{title: "Polls", url: "polls", beta: true},
                    {title: "Predictions", url: "predictions", alpha: true}];
    $scope.location = $location;
})

.controller('HomeController', function($scope, $log) {
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
    $scope.polls = [{"creator": "erb", "text": "This is a poll"},
                    {"creator": "clara", "text": "This is another poll"}];
    $scope.time = new Date().toISOString();
})

.controller('PollController', function($scope, $resource, $log) {
    $scope.voted = true;
    $scope.votes = 0;
    $scope.rating = 0;

    $scope.vote = function(option) {
        $scope.votes += 1;
        $scope.rating = option ? $scope.rating+1 : $scope.rating-1;
        $scope.trues = Math.round(1000*(1+($scope.rating/$scope.votes))/2)/10;
        $scope.falses = Math.round(1000*(1-($scope.rating/$scope.votes))/2)/10;
        console.log($scope.rating);
        $scope.voted = true;
    };
})

.controller('ProfileController', function($scope, $routeParams, User) {
    $scope.user = {};
    User("username", $routeParams.username).$promise.then(function(payload) {
        console.log(payload);
        $scope.user = payload.data[0];
        $scope.email_hash = CryptoJS.MD5($scope.user.email).toString();
    }, function(error) {
        $scope.loading_error = true;
    });
})

.controller('LoginController', function($scope, $routeParams, $location, $http) {
    $scope.logging_in = false;

    $scope.login = function() {
        $scope.logging_in = true;
        console.log("Hello");
        $http.post('/api/0/auth', {username: $scope.username, password: $scope.password})
        .success(function(data, status, headers, config) {
            console.log(data);
            if(!data.auth) {
                $scope.error = data.error;
            } else {
                $scope.error = "";
                $location.path("/profile/"+$scope.username);
                alert("Success!");
            }
            $scope.logging_in = false;
        }).error(function(data, status, headers, config) {
            $scope.error = "Something went wrong, we dearly apologize";
            $scope.logging_in = false;
        });
    };
})

.config(function($routeProvider, $locationProvider) {
  $routeProvider
   .when('/', {
    templateUrl: 'home.html',
    controller: 'HomeController',
  })
   .when('/polls', {
    templateUrl: 'polls.html',
    controller: 'PollsController',
  })
   .when('/profile/:username', {
    templateUrl: 'profile.html',
    controller: 'ProfileController',
  })
   .when('/admin', {
    templateUrl: 'admin.html',
    controller: 'AdminController',
  })
   .when('/login', {
    templateUrl: 'login.html',
    controller: 'LoginController',
  });
});
