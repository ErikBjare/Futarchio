app.controller('MainController', function($scope, $rootScope, $window) {
    $scope.is_logged_in = function() {
        return Meteor.userId() ? true : false;
    };

    $rootScope.$on('$routeChangeSuccess', function(event) {
//        $window.ga('send', 'pageview', { page: $location.path() });
    });
});

app.controller('HomeController', function($scope, $log, $location, $anchorScroll) {
});

app.controller('PollsController', function($scope, $resource, $log) {
});

app.controller('NewPollController', function($scope, $resource, $log, $window) {
    $scope.createPoll = function() {
        var poll = {"title": $scope.title, "description": $scope.description || "", "type": $scope.type};
        console.log(poll);
    };
});

app.controller("StatementController", function($scope, $resource) {
});

app.controller('CreateStatementController', function($scope, $resource, $log, Statement) {
});

app.controller('StatementsController', function($scope, $resource, $log, Statement) {
});

app.controller('PollController', function($scope, $resource, $log) {
});

app.controller('ProfileController', function($scope, $routeParams, $location, $cookieStore, gravatar, User) {
    console.log(Meteor.user().email);
    //$scope.gravatar_hash = gravatar.hash($scope.profile.email);
});
