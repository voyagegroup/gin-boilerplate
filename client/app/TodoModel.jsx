export default class TodoModel {
    constructor() {
        this.onChanges = [];
        this.todos = [];
    }

    subscribe(onChange) {
        this.onChanges.push(onChange);
        this.inform();
    }

    inform() {
        this.onChanges.forEach(cb => cb());
    }

    load() {
        fetch('/api/todos')
        .then(x => x.json())
        .then(json => {
            if (json == null) {
                return;
            }
            this.todos = json;
            this.inform();
        })
        .catch(err => {
            console.error('fetch error', err);
        });
    }

    addTodo(title) {
        const todo = {
            title: title,
            completed: false
        };

        fetch('/api/todos', {
            method: 'PUT',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(todo)
        })
        .then(x => x.json())
        .then(data => {
            this.todos = this.todos.concat([data]);
            this.inform();
        })
        .catch(err => {
            console.error('post todo error: ', err);
        });
    }

    toggleAll(checked) {
        fetch('/api/todos/toggleall', {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({checked: checked})
        })
        .then(resp => this.checkStatus(resp, 200))
        .then(() => {
            this.todos = this.todos.map(todo => {
                return Object.assign({}, todo, {completed: checked});
            });
            this.inform();
        })
        .catch(err => {
            console.error('post todo error: ', err);
        });
    }

    toggle(todoToToggle) {
        fetch('/api/todos/toggle', {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(todoToToggle)
        })
        .then(resp => this.checkStatus(resp, 200))
        .then(() => {
            this.todos = this.todos.map(todo => {
                return todo !== todoToToggle ?
                    todo :
                    Object.assign({}, todo, {completed: !todo.completed});
            });
            this.inform();
        })
        .catch(err => {
            console.error('post todo error: ', err);
        });
    }

    checkStatus(resp, code) {
        if (resp.status == code) {
            return resp
        } else {
            const error = new Error(resp.statusText)
            error.resp = resp
            throw error
        }
    }

    destroy(todo) {
        fetch('/api/todos', {
            method: 'DELETE',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(todo)
        })
        .then(resp => this.checkStatus(resp, 200))
        .then(() => {
            this.todos = this.todos.filter(candidate => {
                return candidate !== todo;
            });
            this.inform();
        })
        .catch(err => {
            console.error('post todo error: ', err);
        });
    }

    save(todoToSave, text) {
        const toSave = Object.assign({}, todoToSave, {title: text});

        fetch('/api/todos', {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(toSave)
        })
        .then(x => x.json())
        .then(() => {
            this.todos = this.todos.map(todo => {
                return todo !== todoToSave ? todo : toSave;
            });
            this.inform();
        })
        .catch(err => {
            console.error('post todo error: ', err);
        });
    }

    clearCompleted() {
        const todosToDelete = this.todos.filter(todo => {
            return todo.completed;
        });

        fetch('/api/todos/multi', {
            method: 'DELETE',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(todosToDelete)
        })
        .then(resp => this.checkStatus(resp, 200))
        .then(() => {
            this.todos = this.todos.filter(todo => {
                return !todo.completed;
            });
            this.inform();
        })
        .catch(err => {
            console.error('delete todo error: ', err);
        });
    }
}

