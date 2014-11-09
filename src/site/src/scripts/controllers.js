app.controller('MainController', function($scope, $route, $rootScope, $location, $window, user) {
    $rootScope.$on('$routeChangeSuccess', function(event) {
        $window.ga('send', 'pageview', { page: $location.path() });
    });

    $scope.links_left = [{title: "Polls", url: "polls", beta: true},
                    {title: "Predictions", url: "predictions", alpha: true}];
    $scope.location = $location;
    $scope.user = user;
});

app.controller('HomeController', function($scope, $log, $location, $anchorScroll) {
    $scope.gotoAbout = function() {
        console.log("asd");
        $location.hash("about");
        $anchorScroll();
        $location.hash("");
    };
});

app.controller('AdminController', function($scope, $resource, msgStack) {
    $scope.resources = ["users", "idklol"];
    $scope.keys = ["email", "id"];
    $scope.stdin = "";

    $scope.log = function(obj) {
        $scope.stdout = msgStack($scope, obj);
        $scope.stdin = "";
    };
});

app.controller('PollsController', function($scope, $resource, $log, Poll) {
    Poll.query(function(data) {
           $scope.polls = data;
    });
});

app.controller('NewPollController', function($scope, $resource, $log, $window, Poll) {
    $scope.createPoll = function() {
        var poll = new Poll(
            {"title": $scope.title,
             "description": $scope.description || "",
             "type": $scope.type});
        console.log(poll);
        poll.$save().then(function() {
            // Dirty way of reloading polls
            $window.location.reload();
        });
    };
});

app.controller('PredictionsController', function($scope, $resource, $log) {
});

app.controller('PollController', function($scope, $resource, $log, Vote) {
    $scope.choices = [{"name": "Choice 1", "value": 0}, {"name": "Choice 2", "value": 0}];
    $scope.$broadcast('refreshSlider');


    if($scope.user.is_logged_in() && $scope.user.has_voted_on($scope.poll.id)) {
        $scope.voted = true;
    } else {
        $scope.voted = false;
    }

    $scope.votes = 0;
    for(var k in $scope.poll.weights) {
        $scope.votes += $scope.poll.weights[k];
    }

    if($scope.poll.type == "YesNoPoll") {
        $scope.rating = $scope.poll.weights.yes - $scope.poll.weights.no;

        $scope.update = function() {
            $scope.trues = Math.round(1000*(1+($scope.rating/$scope.votes))/2)/10;
            $scope.falses = Math.round(1000*(1-($scope.rating/$scope.votes))/2)/10;
        };
        $scope.update();

        $scope.vote = function(option) {
            weights = option ? {"yes": 1} : {"no": 1};
            vote = new Vote({weights: weights});
            vote.$save({pollid: $scope.poll.id}).then(function (data) {
                console.log(data);
                $scope.votes += 1;
                $scope.rating = option ? $scope.rating+1 : $scope.rating-1;
                $scope.update();
                $scope.voted = true;
            }, function(data) {
                $scope.error = data.data.error;
                $log.error(data.data.error);
            });
        };
    } else if($scope.poll.type == "MultichoicePoll") {
        $log.warning("Not implemented");
    } else {
        $log.error("Unknown polltype: " + $scope.poll.type);
    }

    console.log($scope.poll);
});

app.controller('ProfileController', function($scope, $routeParams, $location, $cookieStore, gravatar, UserKeyVal) {
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

    UserKeyVal("username", $routeParams.username).$promise.then(function(payload) {
        console.log(payload);
        $scope.profile = payload[0];
        $scope.profile.gravatar_hash = gravatar.hash($scope.profile.email);
    }, function(error) {
        $scope.loading_error = true;
    });
});

app.controller('LoginController', function($scope, $routeParams, $location, $window, user, $http) {
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
            // Needed in order to ask for remembering username & password
            $window.location.href = "/profile/"+data.username;
            $window.location.reload();
            $scope.logging_in = false;
        }, function(error) {
            $scope.error = error;
            $scope.logging_in = false;
        });
    };
});

app.controller('SignupController', function($scope, $routeParams, $location, User, user) {
    if(user.is_logged_in()) {
        $location.path("/profile/"+user.username());
    }

    $scope.signup = function() {
        $scope.signing_up = true;
        var user = new User({username: $scope.username, password: $scope.password, email: $scope.email});
        user.$save().then(function(data) {
            $scope.error = "";
            $location.path("/signup/success");
            $scope.signing_up = false;
        }, function(data) {
            console.log(data.data);
            $scope.error = data.data.error;
            $scope.signing_up = false;
        });
    };
});

app.controller('LogoutController', function(user) {
    user.logout();
});

