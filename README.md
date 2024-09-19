# Shaffra Skill Assessment

This repository contains solutions to the skill assessment tasks. Each solution is located in its respective folder.

## Prerequisites
- **Docker** (required if you want to use the provided compose files)

## How to Run
1. Navigate to the `task-1` or `task-2` folder, depending on which solution you want to run.
2. Execute the `./start.sh` script, or alternatively, run the Docker commands within the script manually.
3. Once the container is running, execute the necessary Go commands inside the container.
4. (Optional) When you're done, you can stop and remove the Docker containers by running:
   ```bash
   docker compose down -t0
   ```
