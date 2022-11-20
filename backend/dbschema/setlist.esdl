module setlist {
  scalar type Visibility extending enum<Public, Private>;

  type SetList {
    required property title -> str;

    required property genres -> json;

    property text -> json;

    required property visibility -> Visibility;

    # TODO: Add comments 
    # property comments -> Comment;

    link created_by -> user::User;

    multi link songs -> song::Song;

    required property song_order -> json;

    property created_at -> datetime;

    property updated_at -> datetime;
  }
}
