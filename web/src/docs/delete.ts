const dummy = {
    token: "my_dummy_token_again",
    ids: ["my_id_1", "my_id_2"],
};

export default {
    code: {
        "One File": `curl --request DELETE \\\n    --url '${
            import.meta.env.VITE_API
        }/api/upload?token=${
            dummy.token
        }' \\\n    --header 'Content-Type: application/json' \\\n    --data '["${
            dummy.ids[0]
        }"]'`,

        "Multiple Files": `curl --request DELETE \\\n    --url '${
            import.meta.env.VITE_API
        }/api/upload?token=${
            dummy.token
        }' \\\n    --header 'Content-Type: application/json' \\\n    --data '["${
            dummy.ids[0]
        }", "${dummy.ids[1]}", ...]'`,
    },

    response: {
        Error: {
            code: "<error_code>",
            message: "<invalid_file_ids>",
        },

        "200": {
            code: 200,
            message: "<used_storage_bytes>",
        },
    },
};
