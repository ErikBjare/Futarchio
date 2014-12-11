Template.profile.helpers({
    // TODO: Move out of profile helpers, should be used in menubar etc.
    // Also make sure people add public emails to take advantage of it
    gravatarUrl: function(email) {
        return Gravatar.imageUrl(email, {secure: true});
    }
});
