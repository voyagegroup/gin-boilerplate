import * as React from 'react';
import Todo from './Todo';
import TodoModel from './TodoModel';
import TodoFooter from './TodoFooter';

export default class TodoBox extends React.Component {
    constructor(props) {
        super(props);

        this.ALL_TODOS = 'all';
        this.ACTIVE_TODOS = 'active';
        this.COMPLETED_TODOS = 'completed';
        this.ENTER_KEY = 13;

        this.state = {
            nowShowing: props.nowShowing,
            editing: null,
            newTodo: ''
        };
    }

    componentDidMount() {
        this.props.model.load();
        setInterval(() => this.props.model.load(), this.props.pollInterval);
    }

    handleChange(event) {
        this.setState({newTodo: event.target.value});
    }

    handleNewTodoKeyDown(event) {
        if (event.keyCode !== this.ENTER_KEY) {
            return;
        }
        event.preventDefault();
        const val = this.state.newTodo.trim();
        if (val) {
            this.props.model.addTodo(val);
            this.setState({newTodo: ''});
        }
    }

    toggleAll(event) {
        const checked = event.target.checked;
        this.props.model.toggleAll(checked);
    }

    toggle(todoToToggle) {
        this.props.model.toggle(todoToToggle);
    }

    destroy(todo) {
        this.props.model.destroy(todo);
    }

    edit(todo) {
        this.setState({editing: todo.id});
    }

    save(todoToSave, text) {
        this.props.model.save(todoToSave, text);
        this.setState({editing: null});
    }

    cancel() {
        this.setState({editing: null});
    }

    clearCompleted() {
        this.props.model.clearCompleted();
    }

    footerClick(mode) {
        this.setState({nowShowing: mode});
        this.props.model.inform();
    }

    render() {
        const todos = this.props.model.todos;

        const shownTodos = todos.filter(todo => {
            switch (this.state.nowShowing) {
                case this.ACTIVE_TODOS:
                    return !todo.completed;
                case this.COMPLETED_TODOS:
                    return todo.completed;
                default:
                    return true;
            }
        }, this);

        const todoItems = shownTodos.map(todo => {
            return (
                <Todo
                  key={todo.id}
                  todo={todo}
                  onToggle={this.toggle.bind(this, todo)}
                  onDestroy={this.destroy.bind(this, todo)}
                  onEdit={this.edit.bind(this, todo)}
                  editing={this.state.editing === todo.id}
                  onSave={this.save.bind(this, todo)}
                  onCancel={this.cancel}
                />
            );
        }, this);

        const activeTodoCount = todos.reduce((accum, todo) => {
            return todo.completed ? accum : accum + 1;
        }, 0);

        const completedCount = todos.length - activeTodoCount;

        let footer = null;
        if (activeTodoCount || completedCount) {
            footer =
                <TodoFooter
                  count={activeTodoCount}
                  completedCount={completedCount}
                  nowShowing={this.state.nowShowing}
                  onClearCompleted={this.clearCompleted.bind(this)}
                  onFooterClick={this.footerClick.bind(this)}
                />;
        }

        let main = null;
        if (todos.length) {
            main = (
                <section id="main">
                    <input
                      id="toggle-all"
                      type="checkbox"
                      onChange={this.toggleAll.bind(this)}
                      checked={activeTodoCount === 0}
                    />
                    <ul id="todo-list">
                        {todoItems}
                    </ul>
                </section>
            );
        }
        return (
            <div>
                <header id="header">
                    <h1>todos</h1>
                    <input
                      id="new-todo"
                      placeholder="What needs to be done?"
                      value={this.state.newTodo}
                      onChange={this.handleChange.bind(this)}
                      onKeyDown={this.handleNewTodoKeyDown.bind(this)}
                      autoFocus={true}
                    />
                </header>
                {main}
                {footer}
            </div>
        );
    }
}

TodoBox.propTypes = {
    url: React.PropTypes.string,
    pollInterval: React.PropTypes.number,
    nowShowing: React.PropTypes.string,
    model: React.PropTypes.instanceOf(TodoModel)
}

TodoBox.defaultProps = {
    nowShowing: 'all'
}
