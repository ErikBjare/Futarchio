Meteor.startup(function() {
    // code to run on server at startup
    Meteor.publish('polls', function() { return Polls.find(); });
    Meteor.publish('votes', function() { return Votes.find(); });
    Meteor.publish('statements', function() { return Statements.find(); });
    Meteor.publish('predictions', function() { return Predictions.find(); });
});

Votes.deny({
    "insert": function(userId, doc) {
        // TODO: This method for deleting existing votes is nasty,
        // move that logic to votebutton controller.
        existingVote = Votes.findOne({createdBy: userId, post: doc.post, type: "UpDown"});

        if(existingVote === undefined) {
            return false;
        } else {
            // Remove existing vote
            Votes.remove(existingVote._id);

            // Overwrite if vote with differing value already exists
            if(existingVote.value != doc.value) {
                return false;
            }

            // If previous vote was identical in value, don't insert it.
            return true;
        }
    }
});
