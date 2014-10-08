app.controller('MainController', function($scope, $route, $location, user) {
    $scope.links_left = [{title: "Polls", url: "polls", beta: true},
                    {title: "Predictions", url: "predictions", alpha: true}];
    $scope.location = $location;
    $scope.user = user;
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

app.controller('PredictionsController', function($scope, $resource, $log) {
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
});

app.controller('ProfileController', function($scope, $routeParams, $location, $cookieStore, gravatar, User) {
    console.log("Missing username routeParam");
    if(!$routeParams.username) {
        console.log("Missing username routeParam");
        me = $cookieStore.get("me");
        console.log(me);
        if(me !== undefined) {
            console.log("Redirecting to profile");
            $location.path("/profile/"+me.username);
        } else {
            console.log("Redirecting to login");
            $location.path("/login");
        }
        return;
    }

    $scope.user = {};
    User("username", $routeParams.username).$promise.then(function(payload) {
        console.log(payload);
        $scope.profile = payload.data[0];
        $scope.profile.gravatar_hash = gravatar.hash($scope.profile.email);
    }, function(error) {
        $scope.loading_error = true;
    });
});

app.controller('LoginController', function($scope, $routeParams, $location, user, $http) {
    //TODO: Set 30 day cookie expiry upon "Remember me"
    $scope.logging_in = false;

    if(user.is_logged_in()) {
        $location.path("/profile/"+user.username());
    } else {
        $scope.logged_in = false;
    }

    $scope.login = function() {
        $scope.logging_in = true;
        user.login($scope.username, $scope.password).then(function(data) {
            $scope.error = "";
            $location.path("/profile/"+$scope.username);
            $scope.logging_in = false;
        }, function(error) {
            $scope.error = error;
            $scope.logging_in = false;
        });
    };
});

app.controller('LogoutController', function(user) {
    user.logout();
});

