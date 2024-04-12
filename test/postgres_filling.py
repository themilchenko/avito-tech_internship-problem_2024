import random
import requests
import json

features = [i for i in range(1, 1001)]
tags = [i for i in range(1, 1001)]

create_banner_url = "http://127.0.0.1:8080/api/banner"

def create_banner(i, feature, tags):
    banner_content = {
        "title": f"Title for {i}",
        "text": f"Text for {i}",
        "url": f"http://example.com/content/{i}"
    }

    payload = {
        "feature_id": feature,
        "tag_ids": tags,
        "content": banner_content,
        "is_active": True
    }

    response = requests.post(create_banner_url, json=payload, cookies={"token":"5d73f225-aed5-4a73-a7ca-2749266b1ced"})

    if response.status_code == 201:
        print(f"Banner with id {i} was successfully created.")
    else:
        print(f"Creating error for banner {id}: {response.text}")

for i, feature in enumerate(features):
    num_tags = random.randint(1, 5)
    feature_tags = random.sample(tags, num_tags)
    
    create_banner(i, feature, feature_tags)
