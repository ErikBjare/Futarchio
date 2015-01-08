SimpleSchema.debug = true;

var Schemas = {};

Schemas.UserCreated = new SimpleSchema({
    createdBy: {
        type: SimpleSchema.RegEx.Id
    },
    createdAt: {
        type: Date
    }
});

Schemas.Rating = new SimpleSchema([Schemas.UserCreated, {
    score: {
        type: Number,
        decimal: true,
        max: 1,
        min: -1
    }
}]);

Schemas.Post = new SimpleSchema([Schemas.UserCreated, {
    title: {
        type: String
    },
    description: {
        type: String
    },
    // TODO: Remove defaultValue?
    ratings: {
        type: Object,
        blackbox: true,
        defaultValue: {}
    },
    tags: {
        type: [String],
        defaultValue: []
    }
}]);

Schemas.Poll = new SimpleSchema([Schemas.Post, {
    type: {
        type: String,
        allowedValues: ["YesNo", "YesNoRange"]
    }
}]);

Schemas.Statement = new SimpleSchema([Schemas.Post, {}]);

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
        type: SimpleSchema.RegEx.Id
    }
}]);


UserCreated = function() {
    var obj = {};
    obj.createdBy = Meteor.userId();
    obj.createdAt = new Date();

    return obj;
};

Rating = function(score) {
    var obj = UserCreated();
    obj.score = score;

    return obj;
};

Post = function(data) {
    var obj = UserCreated();
    obj.title = data.title;
    obj.description = data.description;
    obj.ratings = {};
    obj.tags = data.tags ? data.tags : [];

    return obj;
};

Polls = new Mongo.Collection("polls");
Polls.attachSchema(Schemas.Poll);
Poll = function(data) {
    var obj = Post(data);
    obj.type = data.type;

    return obj;
};

Statements = new Mongo.Collection("statements");
Statements.attachSchema(Schemas.Statement);
Statement = function(data) {
    var obj = Post(data);

    return obj;
};

Predictions = new Mongo.Collection("predictions");
Predictions.attachSchema(Schemas.Prediction);
Prediction = function(data) {
    var obj = UserCreated();
    obj.credence = Number(data.credence);
    obj.statement = data.statement;

    return obj;
};

Votes = new Mongo.Collection("votes");
Votes.attachSchema(Schemas.Vote);
Vote = function(data) {
    var obj = UserCreated();
    obj.type = data.type;
    obj.value = Number(data.value);
    obj.post = data.post;

    return obj;
};
