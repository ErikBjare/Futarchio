app.controller('MainController', function($scope, $route, $location) {
    $scope.links_left = [{title: "Polls", url: "polls", beta: true},
                    {title: "Predictions", url: "predictions", alpha: true}];
    $scope.location = $location;
});

app.controller('HomeController', function($scope, $log) {
});

app.controller('AdminController', function($scope, $resource, msgStack, User) {
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
});

app.controller('PollsController', function($scope, $resource, $log) {
    $scope.polls = [{"creator": "erb", "text": "This is a poll"},
                    {"creator": "clara", "text": "This is another poll"}];
    $scope.time = new Date().toISOString();
});

app.controller('PollController', function($scope, $resource, $log) {
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

app.controller('ProfileController', function($scope, $routeParams, User) {
    $scope.user = {};
    User("username", $routeParams.username).$promise.then(function(payload) {
        console.log(payload);
        $scope.user = payload.data[0];
        $scope.email_hash = CryptoJS.MD5($scope.user.email).toString();
    }, function(error) {
        $scope.loading_error = true;
    });
});

app.controller('LoginController', function($scope, $routeParams, $location, $http) {
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
});

