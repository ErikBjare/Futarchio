Router.configure({
    layoutTemplate: 'layout',
    loadingTemplate: 'loading',
    notFoundtemplate: 'notFound'
});

// Home
Router.route('/', function() {
    this.render("home");
});

// Votes and Polls
Router.route('/vote', function() {
    this.render("polls");
});

Router.route('/poll/:_id', function () {
    var params = this.params;
    var id = params._id;
});

// Predictions and Statements
Router.route('/predict', function() {
    this.render("statements");
});

Router.route('/statement/:_id', function () {
    var params = this.params;
    var id = params._id;
});

Router.route('/prediction/:_id', function () {
    var params = this.params;
    var id = params._id;
});

// User profiles
Router.route('/profile/:username', function () {
    var params = this.params;
    var username = params.username;
});
