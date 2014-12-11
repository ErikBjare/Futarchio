Schemas = {};

Schemas.UserCreated = new SimpleSchema({
    createdBy: {
        type: String
    },
    createdAt: {
        type: Date
    }
});

Schemas.Post = new SimpleSchema([Schemas.UserCreated, {
    title: {
        type: String
    },
    description: {
        type: String
    }
}]);

Schemas.Poll = new SimpleSchema([Schemas.Post, {
    type: {
        type: String,
        allowedValues: ["YesNo", "YesNoRange"]
    }
}]);

Schemas.Statement = new SimpleSchema([Schemas.Post, {
}]);

Schemas.Prediction = new SimpleSchema([Schemas.UserCreated, {
    credence: {
        type: Number
    },
    statement: {
        type: String
    }
}]);

// TODO: Separate Up/Down votes and votes on polls
Schemas.Vote = new SimpleSchema([Schemas.UserCreated, {
    type: {
        type: String,
        allowedValues: ["UpDown", "YesNo", "YesNoRange"]
    },
    value: {
        type: Number,
        decimal: true
    },
    post: {
        type: String
    }
}]);

function addCreationData(data) {
    data.createdBy = Meteor.userId();
    data.createdAt = new Date();
    return data;
}

Polls = new Mongo.Collection("polls");
Polls.attachSchema(Schemas.Poll);
Poll = function(poll) {
    poll = addCreationData(poll);
    return poll;
};

Statements = new Mongo.Collection("statements");
Statements.attachSchema(Schemas.Statement);
Statement = function(stmt) {
    stmt = addCreationData(stmt);
    return stmt;
};

Predictions = new Mongo.Collection("predictions");
Predictions.attachSchema(Schemas.Prediction);
Prediction = function(pred) {
    pred = addCreationData(pred);
    return pred;
};

Votes = new Mongo.Collection("votes");
Votes.attachSchema(Schemas.Vote);
Vote = function(vote) {
    vote.value = Number(vote.value);
    vote = addCreationData(vote);
    return vote;
};
