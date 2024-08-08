<script lang="ts">
    import { onMount } from "svelte";
    import type { ApiResult, OtomotoData } from "./types";

    let pageData: ApiResult<OtomotoData>;
    let link: string;
    let error: string;
    let submitMessage: string;

    const endpointPrefix = "http://localhost:3000/api/car-listings";

    onMount(async () => {
        try {
            const res = await fetch(endpointPrefix.concat("/otomoto"), {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                },
            });
            pageData = await res.json();
        } catch (e: unknown) {
            console.log(e);
        }
    });

    async function handleSubmit() {
        if (link.length > 0) {
            try {
                const formData = new FormData();
                formData.append("link", link);

                const res = await fetch(endpointPrefix.concat("/link"), {
                    method: "POST",
                    body: formData,
                });

                const data: ApiResult = await res.json();

                if (data.message) {
                    submitMessage = data.message;
                } else if (data.error) {
                    error = data.error;
                }
            } catch (e: unknown) {
                console.log(e);
            }
        }
    }
</script>

<form on:submit|preventDefault={handleSubmit}>
    <input bind:value={link} type="text" name="link" id="link" />
    <button type="submit">submit</button>
    {#if error}
        <p>{error}</p>
    {/if}
</form>
