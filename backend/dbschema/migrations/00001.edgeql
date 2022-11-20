CREATE MIGRATION m1uax3mgynkmnenxmp2anm5altg4xzgafofae53zj3riz72oqvktmq
    ONTO initial
{
  CREATE MODULE setlist IF NOT EXISTS;
  CREATE MODULE song IF NOT EXISTS;
  CREATE MODULE user IF NOT EXISTS;
  CREATE FUTURE nonrecursive_access_policies;
  CREATE SCALAR TYPE setlist::Visibility EXTENDING enum<Public, Private>;
  CREATE TYPE setlist::SetList {
      CREATE PROPERTY created_at -> std::datetime;
      CREATE REQUIRED PROPERTY genres -> std::json;
      CREATE REQUIRED PROPERTY song_order -> std::json;
      CREATE PROPERTY text -> std::json;
      CREATE REQUIRED PROPERTY title -> std::str;
      CREATE PROPERTY updated_at -> std::datetime;
      CREATE REQUIRED PROPERTY visibility -> setlist::Visibility;
  };
  CREATE SCALAR TYPE user::UserType EXTENDING enum<Admin, Looper>;
  CREATE TYPE user::User {
      CREATE MULTI LINK set_lists -> setlist::SetList;
      CREATE PROPERTY created_at -> std::datetime;
      CREATE REQUIRED PROPERTY displayName -> std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(3);
          CREATE CONSTRAINT std::min_len_value(3);
      };
      CREATE REQUIRED PROPERTY email -> std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::min_len_value(8);
      };
      CREATE REQUIRED PROPERTY password -> std::str {
          CREATE CONSTRAINT std::min_len_value(8);
      };
      CREATE PROPERTY social_media -> std::json;
      CREATE PROPERTY updated_at -> std::datetime;
      CREATE REQUIRED PROPERTY user_type -> user::UserType {
          SET default := (user::UserType.Looper);
      };
      CREATE REQUIRED PROPERTY username -> std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(3);
          CREATE CONSTRAINT std::min_len_value(3);
      };
      CREATE PROPERTY website -> std::str;
  };
  ALTER TYPE setlist::SetList {
      CREATE LINK created_by -> user::User;
  };
  CREATE SCALAR TYPE song::LoopingType EXTENDING enum<Repeated, SlightChange, Multipart>;
  CREATE TYPE song::Song {
      CREATE LINK submitter -> user::User;
      CREATE REQUIRED PROPERTY artist -> std::str;
      CREATE PROPERTY bpm -> std::int16;
      CREATE PROPERTY created_at -> std::datetime;
      CREATE REQUIRED PROPERTY genre -> std::str;
      CREATE PROPERTY key -> std::str;
      CREATE PROPERTY layers -> std::json;
      CREATE REQUIRED PROPERTY looping_type -> song::LoopingType;
      CREATE PROPERTY music_link -> std::str;
      CREATE PROPERTY song_link -> std::str;
      CREATE PROPERTY text -> std::json;
      CREATE REQUIRED PROPERTY title -> std::str;
      CREATE PROPERTY updated_at -> std::datetime;
      CREATE PROPERTY video_link -> std::str;
  };
  ALTER TYPE setlist::SetList {
      CREATE MULTI LINK songs -> song::Song;
  };
  ALTER TYPE user::User {
      CREATE MULTI LINK saved_songs -> song::Song;
      CREATE MULTI LINK submitted_songs -> song::Song;
  };
};
