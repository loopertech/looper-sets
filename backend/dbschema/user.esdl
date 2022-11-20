module user {
  scalar type UserType extending enum<Admin, Looper>;

  type User {
    required property email -> str {
      constraint min_len_value(8);
      constraint exclusive;
    }

    required property password -> str {
      constraint min_len_value(8);
    }

    required property username -> str {
      constraint exclusive;
      constraint min_len_value(3);
      constraint max_len_value(3);
    }

    required property displayName -> str {
      constraint exclusive;
      constraint min_len_value(3);
      constraint max_len_value(3);
    }

    required property user_type -> UserType {
      default := UserType.Looper;
    };

    property website -> str;

    property social_media -> json;

    multi link submitted_songs -> song::Song;

    multi link saved_songs -> song::Song;

    multi link set_lists -> setlist::SetList;

    property created_at -> datetime;

    property updated_at -> datetime;
  }

}
