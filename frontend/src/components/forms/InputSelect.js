function InputSelect(props) {
    return (
        <div className="mb-3 ">
            <label htmlFor={props.name} className="form-label">
                {" "}
                {props.title}{" "}
            </label>
            <select
                className={`form-control ${props.className}`}
                id={props.name}
                name={props.name}
                value={props.value}
                onChange={props.handleChange}
                placeholder={props.placeholder}
            >
                <option value="">{props.placeholder}</option>
                {props.options.map(o => (
                    <option
                        className="form-select"
                        key={o.id}
                        value={o.value}
                    >
                        {o.value}
                    </option>
                ))}
            </select>
            <div className={props.errorDiv}>{props.errorMsg}</div>

        </div>
    )
}

export default InputSelect

