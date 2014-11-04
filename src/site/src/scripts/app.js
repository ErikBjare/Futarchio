app = angular.module('FutarchioApp', ["ngResource", "ngAnimate", "ngRoute", "ngCookies", "ngSanitize", "vr.directives.slider"]);

app.config(function($routeProvider, $locationProvider) {
  $routeProvider
   .when('/', {
    templateUrl: '/static/home.html',
    controller: 'HomeController',
  })
   .when('/polls', {
    templateUrl: '/static/polls.html',
    controller: 'PollsController',
  })
   .when('/predictions', {
    templateUrl: '/static/predictions.html',
    controller: 'PredictionsController',
  })
   .when('/profile/:username', {
    templateUrl: '/static/profile.html',
    controller: 'ProfileController',
  })
   .when('/profile', {
    templateUrl: '/static/profile.html',
    controller: 'ProfileController',
  })
   .when('/admin', {
    templateUrl: '/static/admin.html',
    controller: 'AdminController',
  })
   .when('/logout', {
    templateUrl: '/static/logout.html',
    controller: 'LogoutController',
  })
   .when('/login', {
    templateUrl: '/static/login.html',
    controller: 'LoginController',
  })
   .when('/signup', {
    templateUrl: '/static/signup.html',
    controller: 'SignupController',
  });

  $locationProvider.html5Mode(true);
});

app.animation('.slide', function() {
    var NG_HIDE_CLASS = 'ng-hide';
    return {
        beforeAddClass: function(element, className, done) {
            if(className === NG_HIDE_CLASS) {
                element.slideUp(done); 
            }
        },
        removeClass: function(element, className, done) {
            if(className === NG_HIDE_CLASS) {
                element.hide().slideDown(done);
            }
        }
    };
});
