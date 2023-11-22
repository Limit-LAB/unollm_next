from openai import OpenAI

client = OpenAI(
    api_key="chatglm",
    base_url="http://127.0.0.1:11451/v1/chat/completions"
)

stream = client.chat.completions.create(
    model="chatglm_turbo",
    messages=[{"role": "user", "content": "你谁啊"}],
    stream=True,
)
for part in stream:
    print(part.choices[0].delta.content or "")