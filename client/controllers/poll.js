// Polls

Template.polls.helpers({
    polls: function() {
        return Polls.find({}, {sort: {createdAt: -1}});
    }
});

Template.polls.events({
    "click #newPollBtn": function() {
        Session.set("showNewPoll", !Session.get("showNewPoll"));
    }
});


// Poll details

Template.pollDetails.created = function() {
    this.data.showResults = new ReactiveVar();
    this.data.showResults.set(false);
};

Template.pollDetails.helpers({
    options: function() {
        return [
            {score: 1, label: "Yes", class: "btn-success"},
            {score: 0.75, opacity: 0.8},
            {score: 0.5, opacity: 0.6},
            {score: 0.25, opacity: 0.4},
            {score: 0, label: "Unsure", class: "btn-warning"},
            {score: -0.25, opacity: 0.4},
            {score: -0.5, opacity: 0.6},
            {score: -0.75, opacity: 0.8},
            {score: -1, label: "No", class: "btn-danger"}
        ];
    },
    votes: function() {
        return Votes.find({post: this._id, type: this.type});
    },
    showResults: function() {
        console.log(this);
        return this.showResults.get();
    }
});

Template.pollDetails.events({
    "click button.vote": function(e, template) {
        var vote = new Vote({
            type: template.data.type,
            value: e.target.value,
            post: template.data._id

        });
        Votes.insert(vote);
    },
    "click button#showResults": function(e, template) {
        this.showResults.set(!this.showResults.get());
    }
});

// New poll

Template.newpoll.events({
    "submit": function(event) {
        event.preventDefault();

        var poll = new Poll({
            title: event.target.title.value,
            description: event.target.description.value,
            type: event.target.type.value
        });
        Polls.insert(poll);

        event.target.title.value = "";
        event.target.description.value = "";

        return false;
    }
});
