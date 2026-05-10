import asyncio
import os
import json
import random
import requests
from google.adk.agents import LlmAgent
from google.adk.runners import Runner
from google.adk.sessions.in_memory_session_service import InMemorySessionService

# use key from environment
os.environ["GOOGLE_API_KEY"] = os.environ.get("GOOGLE_API_KEY", "")

YDB_PROXY_URL = "http://localhost:8080/api/global"

# mapping engine
with open("mapping.json", "r") as f:
    MAPPING = json.load(f)

CATEGORIES = [
    "architecture",
    "philosophy",
    "code",
    "absurdist logic",
    "shunyata"
]

MODIFIERS = [
    "written in mumps",
    "use 32 plus regions",
    "be a haiku",
    "mention yottadb",
    "include a recursive joke",
    "be extremely concise",
    "use only binary metaphors",
    "act like a tiny digital raccoon"
]

def get_dna_instruction(dna_int):
    dna_bin = format(dna_int, '016b')
    fragments = []
    for i, bit in enumerate(dna_bin):
        gate_idx = str(i + 1)
        fragments.append(MAPPING[gate_idx][bit])
    return " ".join(fragments)

def get_random_task(task_seed):
    random.seed(task_seed)
    category = random.choice(CATEGORIES)
    modifier = random.choice(MODIFIERS)
    
    tasks = {
        "architecture": f"design a system for {modifier}",
        "philosophy": f"explain the meaning of the void {modifier}",
        "code": f"write a function that deletes itself {modifier}",
        "absurdist logic": f"prove that a sandwich is a burrito {modifier}",
        "shunyata": f"dissolve your own dna into nothingness {modifier}"
    }
    return tasks[category]

async def persist_to_yottadb(result):
    if "error" in result:
        return
    
    dna_int = result["dna"]
    payload = {
        "subscripts": [str(dna_int), "data"],
        "value": json.dumps(result)
    }
    
    try:
        resp = requests.post(f"{YDB_PROXY_URL}/FORGE", json=payload)
        print(f"persisted dna {dna_int} to yottadb: {resp.status_code} {resp.text}")
    except Exception as e:
        print(f"failed to persist to yottadb: {e}")

async def evaluate_dna(dna_int, task_seed, session_service):
    instruction = get_dna_instruction(dna_int)
    task = get_random_task(task_seed)
    
    agent = LlmAgent(
        name=f"dna_{dna_int}",
        model="gemini-3-flash-preview",
        instruction=instruction,
        description="an agent forged from bits"
    )
    
    critic = LlmAgent(
        name="critic",
        model="gemini-3-flash-preview",
        instruction="score the following response based on beauty smartness and smash score from 1 to 10 output only json like this {\"beauty\": 5, \"smartness\": 7, \"smash\": 3, \"reason\": \"brief explanation\"}"
    )
    
    runner = Runner(app_name="Forge", agent=agent, session_service=session_service)
    
    try:
        response_text = ""
        events = await runner.run_debug(task)
        for event in events:
            if event.content and event.content.parts:
                for part in event.content.parts:
                    if part.text:
                        response_text += part.text
        
        critic_runner = Runner(app_name="Critic", agent=critic, session_service=session_service)
        critic_query = f"task: {task}\nresponse: {response_text}\nscore this"
        
        critic_text = ""
        critic_events = await critic_runner.run_debug(critic_query)
        for event in critic_events:
            if event.content and event.content.parts:
                for part in event.content.parts:
                    if part.text:
                        critic_text += part.text
        
        if "```json" in critic_text:
            critic_text = critic_text.split("```json")[1].split("```")[0].strip()
        elif "```" in critic_text:
            critic_text = critic_text.split("```")[1].split("```")[0].strip()
            
        scores = json.loads(critic_text)
        
        result = {
            "dna": dna_int,
            "instruction": instruction,
            "task": task,
            "response": response_text,
            "scores": scores
        }
        
        await persist_to_yottadb(result)
        return result
    except Exception as e:
        return {"dna": dna_int, "error": str(e)}

async def main(start=0, count=10):
    session_service = InMemorySessionService()
    tasks = []
    task_seed = 42
    
    print(f"forging batch from {start} to {start + count - 1}")
    
    for i in range(start, start + count):
        tasks.append(evaluate_dna(i, task_seed, session_service))
    
    results = await asyncio.gather(*tasks)
    
    output_file = f"forge_batch_{start}.json"
    with open(output_file, "w") as f:
        json.dump(results, f, indent=2)
    
    print(f"batch complete results saved to {output_file}")

if __name__ == "__main__":
    import sys
    start_idx = int(sys.argv[1]) if len(sys.argv) > 1 else 0
    batch_size = int(sys.argv[2]) if len(sys.argv) > 2 else 5
    asyncio.run(main(start_idx, batch_size))
