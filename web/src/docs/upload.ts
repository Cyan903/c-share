const dummy = {
    token: "my_dummy_token",
    comment: "file_comment",
    password: "my_password_123",
    path: "/path/to/file.png",
};

export default {
    code: {
        Public: `curl --request POST \\\n    --url '${
            import.meta.env.VITE_API
        }/api/upload?token=${dummy.token}&perm=public&comment=${
            dummy.comment
        }' \\\n    --header 'Content-Type: multipart/form-data' \\\n    --form upload=@${
            dummy.path
        }`,

        Private: `curl --request POST \\\n    --url '${
            import.meta.env.VITE_API
        }/api/upload?token=${dummy.token}&perm=private&comment=${
            dummy.comment
        }' \\\n    --header 'Content-Type: multipart/form-data' \\\n    --form upload=@${
            dummy.path
        }`,

        Unlisted: `curl --request POST \\\n    --url '${
            import.meta.env.VITE_API
        }/api/upload?token=${dummy.token}&perm=unlisted&password=${
            dummy.password
        }&comment=${
            dummy.comment
        }' \\\n    --header 'Content-Type: multipart/form-data' \\\n    --form upload=@${
            dummy.path
        }`,
    },

    response: {
        Error: {
            code: "<error_code>",
            message: "<error_message>",
        },

        "200": {
            code: 200,
            count: 2,
            data: {
                id: "<file_id>",
                storage: "<used_storage_bytes>",
            },
        },
    },
};
