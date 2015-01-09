// Polls

var defaultOrder = "date";


Template.polls.created = function() {
    TemplateVar.set(this, "order", defaultOrder);
};

Template.polls.helpers({
    polls: function(order) {
        // TODO: Sort by points & activity by publishing and subscribing to aggregations server-side
        // https://stackoverflow.com/questions/18520567/average-aggregation-queries-in-meteor
        // Votes.aggregate([{$match: {type: "UpDown"}}, {$group: {_id: "$post", score: {$sum: "$value"}}}]);

        var orderToSort = {
            "date": {createdAt: -1},
            "points": {points: -1},
            "activity": {createdAt: -1}
        };
        polls = Polls.find({}, {sort: orderToSort[order]}).fetch();
        polls = _.map(polls, function(obj) {
            return _.extend(obj, {"cardType": "poll"});
        });
        return polls;
    },
    order: function() {
        return TemplateVar.get("order");
    }
});

Template.polls.events({
    "click #orderMenu": function(e, template) {
        if(e.target.id) {
            TemplateVar.set(template, "order", e.target.id);
        }
    },
    "click #newPollBtn": function(e, template) {
        TemplateVar.set(template, "showNewPoll", !TemplateVar.get(template, "showNewPoll"));
    }
});


// Poll details

Template.pollDetails.created = function() {
    TemplateVar.set(this, "showResults", false);
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
        return TemplateVar.get("showResults");
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
        TemplateVar.set(template, "showResults", !TemplateVar.get(template, "showResults"));
    }
});

// New poll

Template.newpoll.events({
    "submit": function(event) {
        event.preventDefault();

        // Empty variable field handling
        if (!event.target.title.value){
            Session.set("formerror", "You cannot create a question without a title");
            return false;
        }
        else if (!event.target.description.value){
            Session.set("formerror", "You cannot create a question without a description");
            return false;
        }

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
