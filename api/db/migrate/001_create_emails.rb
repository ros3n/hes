class CreateEmails < ActiveRecord::Migration[5.2]
  def up
    create_table :emails do |t|
      t.string :user_id
      t.string :sender
      t.string :recipients, array: true, default: []
      t.string :subject
      t.string :message
      t.string :status
    end
  end

  def down
    drop_table :emails
  end
end
