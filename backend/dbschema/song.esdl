module song {
  scalar type LoopingType extending enum<Repeated, SlightChange, Multipart>;

  type Song {
    required property title -> str;

    required property artist -> str;

    required property genre -> str;

    required property looping_type -> LoopingType;

    property bpm -> int16;

    property key -> str;

    property layers -> json;

    property text -> json;

    property video_link -> str;

    property song_link -> str;

    property music_link -> str;

    # TODO: Add comments 
    # property comments -> Comment;

    link submitter -> user::User;
    
    property created_at -> datetime;

    property updated_at -> datetime;
  }
}
